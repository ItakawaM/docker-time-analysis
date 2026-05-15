package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ItakawaM/docker-time-analysis/internal/compute"
	"github.com/ItakawaM/docker-time-analysis/internal/parse"
)

func (s *Server) writeJSONResponse(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func (s *Server) validateRequest(w http.ResponseWriter, r *http.Request, method string, contentType string) bool {
	if r.Method != method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), contentType) {
		http.Error(w, fmt.Sprintf("Expected %s data", contentType), http.StatusBadRequest)
		return false
	}

	return true
}

func (s *Server) HandleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 100<<20) // 100mb limit
	if !s.validateRequest(w, r, http.MethodPost, FormData) {
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

	if err := s.writeJSONResponse(w, http.StatusOK, UploadResponse{
		Message:    "Parsed the file",
		ParsedRows: len(parsedData),
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *Server) HandleCompute(w http.ResponseWriter, r *http.Request) {
	if !s.validateRequest(w, r, http.MethodPost, JSON) {
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
		http.Error(w, "No data loaded, call /upload first", http.StatusBadRequest)
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

	alpha, beta := table.ComputeLinearRegressionParams()
	rSquared := table.ComputeRSquared(alpha, beta)

	s.mu.Lock()
	s.alpha = alpha
	s.beta = beta
	s.table = table
	s.mu.Unlock()

	if err := s.writeJSONResponse(w, http.StatusOK, NewComputeResponse(sample, table, alpha, beta, rSquared)); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *Server) HandleSignificance(w http.ResponseWriter, r *http.Request) {
	if !s.validateRequest(w, r, http.MethodPost, JSON) {
		return
	}

	var request SignificanceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Failed to parse json data", http.StatusBadRequest)
		return
	}

	if request.SignificanceLevel <= 0 || request.SignificanceLevel >= 1 {
		http.Error(w, "Significance level must be between 0 and 1", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	table, alpha, beta := s.table, s.alpha, s.beta
	s.mu.RUnlock()

	if table == nil {
		http.Error(w, "No computation loaded, call /compute first", http.StatusBadRequest)
		return
	}

	fisherEmpirical, fisherCritical := table.ComputeFisherStatistics(alpha, beta, request.SignificanceLevel, 2)
	pearsonCorrelation, pearsonEmpirical, pearsonCritical := table.ComputePearsonCorrelation(request.SignificanceLevel)

	if err := s.writeJSONResponse(w, http.StatusOK, NewSignificanceResponse(fisherEmpirical, fisherCritical,
		pearsonCorrelation, pearsonEmpirical, pearsonCritical)); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
