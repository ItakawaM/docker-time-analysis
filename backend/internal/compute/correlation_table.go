package compute

import (
	"fmt"
	"math"

	"github.com/ItakawaM/docker-time-analysis/internal/parse"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

type interval struct {
	begin float64
	end   float64
	mid   float64
}

func buildIntervals(min float64, max float64, count int) ([]interval, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count has to be positive, got = %d", count)
	}

	width := (max - min) / float64(count)
	intervals := make([]interval, count)
	for i := range intervals {
		begin := min + float64(i)*width
		end := begin + width

		intervals[i] = interval{
			begin: begin,
			end:   end,
			mid:   (begin + end) / 2,
		}
	}

	return intervals, nil
}

type CorrelationTable struct {
	totalValues int

	yIntervals []interval
	xIntervals []interval

	// (y, x) -> row: y, column: x
	frequencies *mat.Dense

	xMarginals       *mat.VecDense // column totals
	yMarginals       *mat.VecDense // row totals
	conditionalMeanY *mat.VecDense
}

func NewCorrelationTable(data []*parse.DockerEntry) (*CorrelationTable, error) {
	minX, maxX := data[0].DockerCount, data[0].DockerCount
	minY, maxY := data[0].StartupTime, data[0].StartupTime

	for _, entry := range data {
		x, y := entry.DockerCount, entry.StartupTime
		if x < minX {
			minX = x
		} else if x > maxX {
			maxX = x
		}

		if y < minY {
			minY = y
		} else if y > maxY {
			maxY = y
		}
	}

	intervals, _ := SturgesCoeff(len(data))
	xIntervals, err := buildIntervals(minX, maxX, intervals)
	if err != nil {
		return nil, err
	}

	yIntervals, err := buildIntervals(minY, maxY, intervals)
	if err != nil {
		return nil, err
	}

	freq := mat.NewDense(intervals, intervals, nil)
	table := &CorrelationTable{
		totalValues:      len(data),
		xIntervals:       xIntervals,
		yIntervals:       yIntervals,
		frequencies:      freq,
		xMarginals:       mat.NewVecDense(len(xIntervals), nil),
		yMarginals:       mat.NewVecDense(len(yIntervals), nil),
		conditionalMeanY: mat.NewVecDense(len(xIntervals), nil),
	}

	for _, entry := range data {
		// get (i, j) -> row: i, col: j
		i := table.findInterval(yIntervals, entry.StartupTime)
		j := table.findInterval(xIntervals, entry.DockerCount)
		if i >= 0 && j >= 0 {
			table.frequencies.Set(i, j, table.frequencies.At(i, j)+1)
		}
	}

	table.assignMarginals()
	table.assignConditionalMean()

	return table, nil
}

func (ct *CorrelationTable) findInterval(intervals []interval, value float64) int {
	for i, interval := range intervals {
		if value >= interval.begin && (value < interval.end || i == len(intervals)-1) {
			return i
		}
	}

	return -1
}

func (ct *CorrelationTable) assignMarginals() {
	identityVector := mat.NewVecDense(len(ct.yIntervals), nil)
	for k := range ct.yIntervals {
		identityVector.SetVec(k, 1)
	}
	// Matrix[NxM] -> Matrix[MxN] | Matrix[MxN] x Identity[Nx1] -> ColumnSumVector[Nx1]
	ct.xMarginals.MulVec(ct.frequencies.T(), identityVector)

	identityVector = mat.NewVecDense(len(ct.xIntervals), nil)
	for k := range ct.xIntervals {
		identityVector.SetVec(k, 1)
	}
	// Matrix[NxM] x Identity[Mx1] -> RowSumVector[Mx1]
	ct.yMarginals.MulVec(ct.frequencies, identityVector)
}

func (ct *CorrelationTable) assignConditionalMean() {
	yMids := mat.NewVecDense(len(ct.yIntervals), nil)
	for i, interval := range ct.yIntervals {
		yMids.SetVec(i, interval.mid)
	}
	ct.conditionalMeanY.MulVec(ct.frequencies.T(), yMids)
	for i := range ct.xIntervals {
		// Sum(y_i * n_ij) / n_i
		ct.conditionalMeanY.SetVec(i,
			ct.conditionalMeanY.AtVec(i)/ct.xMarginals.AtVec(i))
	}
}

func (ct *CorrelationTable) ComputeLinearRegressionParams() (alpha float64, beta float64) {
	x := ct.GetXMids()
	y := ct.GetConditionalMeanY()
	weights := ct.GetXMarginals()

	// y = ax + b
	beta, alpha = stat.LinearRegression(x, y, weights, false)
	return
}

func (ct *CorrelationTable) ComputeRSquared(alpha float64, beta float64) (rSquared float64) {
	Q, _, Qo := ct.computeVariations(alpha, beta)
	rSquared = 1 - Qo/Q
	return
}

func (ct *CorrelationTable) ComputeFisherStatistics(alpha float64, beta float64, significanceLevel float64, mParams int) (empirical float64, critical float64) {
	_, Qp, Qo := ct.computeVariations(alpha, beta)
	df1, df2 := float64(mParams-1), float64(mParams)
	empirical = Qp * df2 / (Qo * df1)

	critical = distuv.F{
		D1: df1,
		D2: df2,
	}.Quantile(1 - significanceLevel)

	return
}

func (ct *CorrelationTable) computeVariations(alpha, beta float64) (Q, Qp, Qo float64) {
	x := ct.GetXMids()
	y := ct.GetConditionalMeanY()
	weights := ct.GetXMarginals()
	yMean := stat.Mean(y, weights)

	for i := range x {
		ni := weights[i]
		yTheoretical := beta + alpha*x[i]

		diffTotal := y[i] - yMean
		diffResidual := y[i] - yTheoretical
		diffRegression := yTheoretical - yMean

		Q += ni * diffTotal * diffTotal            // total variation
		Qo += ni * diffResidual * diffResidual     // residual variation
		Qp += ni * diffRegression * diffRegression // regression variation
	}
	return
}

func (ct *CorrelationTable) ComputePearsonCorrelation(significanceLevel float64) (r float64, empirical float64, critical float64) {
	x := ct.GetXMids()
	y := ct.GetConditionalMeanY()
	weights := ct.GetXMarginals()

	r = stat.Correlation(x, y, weights)

	df := float64(ct.totalValues - 2.0)
	empirical = r * math.Sqrt(df) / math.Sqrt(1.0-r*r)

	distribution := distuv.StudentsT{
		Mu:    0,
		Sigma: 1,
		Nu:    df,
	}
	critical = distribution.Quantile(1.0 - significanceLevel/2)

	return
}

func (ct *CorrelationTable) GetXMids() []float64 {
	xMids := make([]float64, len(ct.xIntervals))
	for i := range xMids {
		xMids[i] = ct.xIntervals[i].mid
	}

	return xMids
}

func (ct *CorrelationTable) GetYMids() []float64 {
	yMids := make([]float64, len(ct.yIntervals))
	for i := range yMids {
		yMids[i] = ct.yIntervals[i].mid
	}

	return yMids
}

func (ct *CorrelationTable) GetXMarginals() []float64 {
	xMarginals := make([]float64, len(ct.xMarginals.RawVector().Data))
	for i := range xMarginals {
		xMarginals[i] = ct.xMarginals.AtVec(i)
	}

	return xMarginals
}

func (ct *CorrelationTable) GetYMarginals() []float64 {
	yMarginals := make([]float64, len(ct.yMarginals.RawVector().Data))
	for i := range yMarginals {
		yMarginals[i] = ct.yMarginals.AtVec(i)
	}

	return yMarginals
}

func (ct *CorrelationTable) GetConditionalMeanY() []float64 {
	conditionalY := make([]float64, len(ct.conditionalMeanY.RawVector().Data))
	for i := range conditionalY {
		conditionalY[i] = ct.conditionalMeanY.AtVec(i)
	}

	return conditionalY
}

func (ct *CorrelationTable) GetFrequencies() []float64 {
	return append([]float64(nil), ct.frequencies.RawMatrix().Data...)
}
