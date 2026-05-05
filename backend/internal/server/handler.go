package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ItakawaM/docker-time-analysis/internal/compute"
	"github.com/ItakawaM/docker-time-analysis/internal/parse"
	"gonum.org/v1/gonum/stat"
)

func (s *Server) HandleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 100<<20) // 100mb limit

	// Redundant check, but better be sure
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		http.Error(w, "Expected multipart form data", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(100 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	parsedData, err := parse.LoadDataFromFormFile(file)
	if err != nil {
		http.Error(w, "Invalid .csv file provided", http.StatusUnprocessableEntity)
		return
	}

	if len(parsedData) < 50 {
		http.Error(w, "The .csv file has to contain at least 50 entries", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.data = parsedData
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(UploadResponse{
		Message:    "Parsed the file",
		ParsedRows: len(parsedData),
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *Server) HandleCompute(w http.ResponseWriter, r *http.Request) {
	// Redundant check, but better be sure
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "Expected json data", http.StatusBadRequest)
		return
	}

	var request ComputeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Failed to parse json data", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	data := s.data
	s.mu.RUnlock()

	if len(data) == 0 {
		http.Error(w, "No data loaded", http.StatusBadRequest)
		return
	}

	if request.SampleSize < 50 || request.SampleSize > len(data) {
		http.Error(w, fmt.Sprintf("Sample size is out of bounds: %d", request.SampleSize), http.StatusBadRequest)
		return
	}

	sample, err := compute.GetSample(data, request.SampleSize)
	if err != nil {
		http.Error(w, "Failed to get sample", http.StatusBadRequest)
		return
	}

	table, err := compute.NewCorrelationTable(sample)
	if err != nil {
		http.Error(w, "Incorrect data provideed", http.StatusBadRequest)
		return
	}

	xMids := table.GetXMids()
	weights := make([]float64, len(xMids))
	for i := range weights {
		weights[i] = table.XMarginal.AtVec(i)
	}

	beta, alpha := stat.LinearRegression(xMids, table.GetYMids(), weights, false)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(NewComputeResponse(sample, table, alpha, beta)); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
