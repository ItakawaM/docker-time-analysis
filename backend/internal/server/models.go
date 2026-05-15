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

	AlphaCoefficient float64 `json:"alphaCoefficient"`
	BetaCoefficient  float64 `json:"betaCoefficient"`

	RSquared float64 `json:"rSquared"`
}

type CorrelationTableData struct {
	YMids []float64 `json:"yMids"`
	XMids []float64 `json:"xMids"`

	Frequencies []float64 `json:"frequencies"`

	XMarginal []float64 `json:"xMarginal"`
	YMarginal []float64 `json:"yMarginal"`

	ConditionalMeanY []float64 `json:"conditionalMeanY"`
}

func NewComputeResponse(sample []*parse.DockerEntry, table *compute.CorrelationTable, alpha float64, beta float64, rSquared float64) *ComputeResponse {
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
			YPoints:          yPoints,
			XPoints:          xPoints,
			AlphaCoefficient: alpha,
			BetaCoefficient:  beta,
			RSquared:         rSquared,
		},
	}
}
