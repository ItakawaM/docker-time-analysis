package server

import (
	"github.com/ItakawaM/docker-time-analysis/internal/compute"
	"github.com/ItakawaM/docker-time-analysis/internal/parse"
)

// UploadResponse contains the result of a successful CSV file upload.
type UploadResponse struct {
	Message    string `json:"message"`
	ParsedRows int    `json:"parsedRows"`
}

// ComputeRequest specifies parameters for regression computation, including the sample size to use from uploaded data.
type ComputeRequest struct {
	SampleSize int `json:"sampleSize"`
}

// ComputeResponse contains the complete results from regression computation including correlation table data and regression models.
type ComputeResponse struct {
	CorrelationTableData CorrelationTableData `json:"correlationTableData"`
	RegressionData       RegressionData       `json:"regressionData"`
}

// RegressionData contains scatter plot points and computed regression models.
type RegressionData struct {
	YPoints []float64 `json:"yPoints"`
	XPoints []float64 `json:"xPoints"`

	LinearRegression      compute.LinearRegression      `json:"linearRegression"`
	PiecewiseRegression   compute.PiecewiseRegression   `json:"piecewiseRegression"`
	ExponentialRegression compute.ExponentialRegression `json:"exponentialRegression"`
}

// CorrelationTableData contains frequency distribution and statistical data for visualization and analysis.
type CorrelationTableData struct {
	YMids []float64 `json:"yMids"`
	XMids []float64 `json:"xMids"`

	Frequencies []float64 `json:"frequencies"`

	XMarginal []float64 `json:"xMarginal"`
	YMarginal []float64 `json:"yMarginal"`

	ConditionalMeanY []float64 `json:"conditionalMeanY"`
}

// SignificanceRequest specifies the significance level for statistical tests.
type SignificanceRequest struct {
	SignificanceLevel float64 `json:"significanceLevel"`
}

// SignificanceResponse contains statistical test results for model adequacy and correlation significance.
type SignificanceResponse struct {
	FisherLinear      compute.StatTestResult `json:"fisherLinear"`
	FisherPiecewise   compute.StatTestResult `json:"fisherPiecewise"`
	FisherExponential compute.StatTestResult `json:"fisherExponential"`
	Pearson           compute.StatTestResult `json:"pearson"`
}

// NewComputeResponse creates a ComputeResponse from sample data, correlation table, and regression models.
func NewComputeResponse(sample []*parse.DockerEntry, table *compute.CorrelationTable,
	linearRegression compute.LinearRegression, piecewiseRegression compute.PiecewiseRegression,
	exponentialRegression compute.ExponentialRegression,
) *ComputeResponse {
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
			YPoints:               yPoints,
			XPoints:               xPoints,
			LinearRegression:      linearRegression,
			PiecewiseRegression:   piecewiseRegression,
			ExponentialRegression: exponentialRegression,
		},
	}
}
