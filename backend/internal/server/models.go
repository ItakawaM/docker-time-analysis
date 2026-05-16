package server

import (
	"github.com/ItakawaM/docker-time-analysis/internal/compute"
	"github.com/ItakawaM/docker-time-analysis/internal/parse"
)

type UploadResponse struct {
	Message    string `json:"message"`
	ParsedRows int    `json:"parsedRows"`
}

type ComputeRequest struct {
	SampleSize int `json:"sampleSize"`
}

type ComputeResponse struct {
	CorrelationTableData CorrelationTableData `json:"correlationTableData"`
	RegressionData       RegressionData       `json:"regressionData"`
}

type RegressionData struct {
	YPoints []float64 `json:"yPoints"`
	XPoints []float64 `json:"xPoints"`

	LinearRegression    compute.LinearRegression    `json:"linearRegression"`
	PiecewiseRegression compute.PiecewiseRegression `json:"piecewiseRegression"`
}

type CorrelationTableData struct {
	YMids []float64 `json:"yMids"`
	XMids []float64 `json:"xMids"`

	Frequencies []float64 `json:"frequencies"`

	XMarginal []float64 `json:"xMarginal"`
	YMarginal []float64 `json:"yMarginal"`

	ConditionalMeanY []float64 `json:"conditionalMeanY"`
}

type SignificanceRequest struct {
	SignificanceLevel float64 `json:"significanceLevel"`
}

type SignificanceResponse struct {
	FisherLinear    compute.StatTestResult `json:"fisherLinear"`
	FisherPiecewise compute.StatTestResult `json:"fisherPiecewise"`
	Pearson         compute.StatTestResult `json:"pearson"`
}

func NewComputeResponse(sample []*parse.DockerEntry, table *compute.CorrelationTable,
	linearRegression compute.LinearRegression, piecewiseRegression compute.PiecewiseRegression) *ComputeResponse {
	yPoints, xPoints := make([]float64, len(sample)), make([]float64, len(sample))
	for i := range sample {
		yPoints[i] = sample[i].StartupTime
		xPoints[i] = sample[i].DockerCount
	}

	return &ComputeResponse{
		CorrelationTableData: CorrelationTableData{
			YMids:            table.GetYMids(),
			XMids:            table.GetXMids(),
			Frequencies:      table.GetFrequencies(),
			XMarginal:        table.GetXMarginals(),
			YMarginal:        table.GetYMarginals(),
			ConditionalMeanY: table.GetConditionalMeanY(),
		},
		RegressionData: RegressionData{
			YPoints:             yPoints,
			XPoints:             xPoints,
			LinearRegression:    linearRegression,
			PiecewiseRegression: piecewiseRegression,
		},
	}
}

func NewSignificanceResponse(fisherLinear compute.StatTestResult, fisherPiecewise compute.StatTestResult, pearson compute.StatTestResult) *SignificanceResponse {
	return &SignificanceResponse{
		FisherLinear:    fisherLinear,
		FisherPiecewise: fisherPiecewise,
		Pearson:         pearson,
	}
}
