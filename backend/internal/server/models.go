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

type GraphData struct {
	YPoints          []float64 `json:"yPoints"`
	XPoints          []float64 `json:"xPoints"`
	AlphaCoefficient float64   `json:"alphaCoefficient"`
	BetaCoefficient  float64   `json:"betaCoefficient"`
}

type TableData struct {
	YMids []float64 `json:"yMids"`
	XMids []float64 `json:"xMids"`

	Frequencies []float64 `json:"frequencies"`

	XMarginal []float64 `json:"xMarginal"`
	YMarginal []float64 `json:"yMarginal"`

	ConditionalMeanY []float64 `json:"conditionalMeanY"`
}

type ComputeResponse struct {
	TableData TableData `json:"tableData"`
	GraphData GraphData `json:"graphData"`
}

func NewComputeResponse(sample []*parse.DockerEntry, table *compute.CorrelationTable, alpha float64, beta float64) *ComputeResponse {
	yPoints, xPoints := make([]float64, len(sample)), make([]float64, len(sample))
	for i := range sample {
		yPoints[i] = sample[i].StartupTime
		xPoints[i] = sample[i].DockerCount
	}

	return &ComputeResponse{
		GraphData: GraphData{
			YPoints:          yPoints,
			XPoints:          xPoints,
			AlphaCoefficient: alpha,
			BetaCoefficient:  beta,
		},
		TableData: TableData{
			YMids: table.GetYMids(),
			XMids: table.GetXMids(),
			// Safe copy
			Frequencies:      append([]float64(nil), table.Frequencies.RawMatrix().Data...),
			XMarginal:        append([]float64(nil), table.XMarginal.RawVector().Data...),
			YMarginal:        append([]float64(nil), table.YMarginal.RawVector().Data...),
			ConditionalMeanY: append([]float64(nil), table.ConditionalMeanY.RawVector().Data...),
		},
	}
}
