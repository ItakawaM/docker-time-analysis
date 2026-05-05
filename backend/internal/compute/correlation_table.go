package compute

import (
	"fmt"

	"github.com/ItakawaM/docker-time-analysis/internal/parse"
	"gonum.org/v1/gonum/mat"
)

type Interval struct {
	Begin float64
	End   float64
	Mid   float64
}

func BuildIntervals(min float64, max float64, count int) ([]Interval, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count has to be positive, got = %d", count)
	}

	width := (max - min) / float64(count)
	intervals := make([]Interval, count)
	for i := range intervals {
		begin := min + float64(i)*width
		end := begin + width

		intervals[i] = Interval{
			Begin: begin,
			End:   end,
			Mid:   (begin + end) / 2,
		}
	}

	return intervals, nil
}

type CorrelationTable struct {
	YIntervals []Interval
	XIntervals []Interval

	// (y, x) -> row: y, column: x
	Frequencies *mat.Dense

	XMarginal        *mat.VecDense // column totals
	YMarginal        *mat.VecDense // row totals
	ConditionalMeanY *mat.VecDense
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
	xIntervals, err := BuildIntervals(minX, maxX, intervals)
	if err != nil {
		return nil, err
	}

	yIntervals, err := BuildIntervals(minY, maxY, intervals)
	if err != nil {
		return nil, err
	}

	freq := mat.NewDense(intervals, intervals, nil)
	table := &CorrelationTable{
		XIntervals:       xIntervals,
		YIntervals:       yIntervals,
		Frequencies:      freq,
		XMarginal:        mat.NewVecDense(len(xIntervals), nil),
		YMarginal:        mat.NewVecDense(len(yIntervals), nil),
		ConditionalMeanY: mat.NewVecDense(len(xIntervals), nil),
	}

	for _, entry := range data {
		// get (i, j) -> row: i, col: j
		i := table.findInterval(yIntervals, entry.StartupTime)
		j := table.findInterval(xIntervals, entry.DockerCount)
		if i >= 0 && j >= 0 {
			table.Frequencies.Set(i, j, table.Frequencies.At(i, j)+1)
		}
	}

	table.calculateMarginals()
	table.calculateConditionalMean()

	return table, nil
}

func (ct *CorrelationTable) findInterval(intervals []Interval, value float64) int {
	for i, interval := range intervals {
		if value >= interval.Begin && (value < interval.End || i == len(intervals)-1) {
			return i
		}
	}

	return -1
}

func (ct *CorrelationTable) calculateMarginals() {
	// Transposition and Multiply by identity vector trick
	identityVector := mat.NewVecDense(len(ct.YIntervals), nil)
	for k := range ct.YIntervals {
		identityVector.SetVec(k, 1)
	}
	// Matrix[NxM] -> Matrix[MxN] | Matrix[MxN] x Identity[Nx1] -> ColumnSumVector[Nx1]
	ct.XMarginal.MulVec(ct.Frequencies.T(), identityVector)

	identityVector = mat.NewVecDense(len(ct.XIntervals), nil)
	for k := range ct.XIntervals {
		identityVector.SetVec(k, 1)
	}
	// Matrix[NxM] x Identity[Mx1] -> RowSumVector[Mx1]
	ct.YMarginal.MulVec(ct.Frequencies, identityVector)
}

func (ct *CorrelationTable) calculateConditionalMean() {
	yMids := mat.NewVecDense(len(ct.YIntervals), nil)
	for i, interval := range ct.YIntervals {
		yMids.SetVec(i, interval.Mid)
	}
	ct.ConditionalMeanY.MulVec(ct.Frequencies.T(), yMids)
	for i := range ct.XIntervals {
		// Sum(y_i * n_ij) / n_i
		ct.ConditionalMeanY.SetVec(i,
			ct.ConditionalMeanY.AtVec(i)/ct.XMarginal.AtVec(i))
	}
}

func (ct *CorrelationTable) GetXMids() []float64 {
	xMids := make([]float64, len(ct.XIntervals))
	for i := range xMids {
		xMids[i] = ct.XIntervals[i].Mid
	}

	return xMids
}

func (ct *CorrelationTable) GetYMids() []float64 {
	yMids := make([]float64, len(ct.YIntervals))
	for i := range yMids {
		yMids[i] = ct.YIntervals[i].Mid
	}

	return yMids
}

func (ct *CorrelationTable) GetConditionalMeanY() []float64 {
	conditionalY := make([]float64, len(ct.ConditionalMeanY.RawVector().Data))
	for i := range conditionalY {
		conditionalY[i] = ct.ConditionalMeanY.AtVec(i)
	}

	return conditionalY
}
