package server

import "github.com/ItakawaM/docker-time-analysis/internal/compute"

type UploadResponse struct {
	Message    string `json:"message"`
	ParsedRows int    `json:"parsedRows"`
}

type ComputeRequest struct {
	SampleSize int `json:"sampleSize"`
}

type ComputeResponse struct {
	YMids []float64 `json:"yMids"`
	XMids []float64 `json:"xMids"`

	Frequencies []float64 `json:"frequencies"`

	XMarginal []float64 `json:"xMarginal"`
	YMarginal []float64 `json:"yMarginal"`

	ConditionalMeanY []float64 `json:"conditionalMeanY"`
}

func NewComputeResponse(table *compute.CorrelationTable) *ComputeResponse {
	yMids := make([]float64, len(table.YIntervals))
	for i := range yMids {
		yMids[i] = table.YIntervals[i].Mid
	}

	xMids := make([]float64, len(table.XIntervals))
	for i := range xMids {
		xMids[i] = table.XIntervals[i].Mid
	}

	return &ComputeResponse{
		YMids: yMids,
		XMids: xMids,
		// Safe copy
		Frequencies:      append([]float64(nil), table.Frequencies.RawMatrix().Data...),
		XMarginal:        append([]float64(nil), table.XMarginal.RawVector().Data...),
		YMarginal:        append([]float64(nil), table.YMarginal.RawVector().Data...),
		ConditionalMeanY: append([]float64(nil), table.ConditionalMeanY.RawVector().Data...),
	}
}
