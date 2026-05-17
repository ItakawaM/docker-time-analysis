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

// HandleUpload handles the POST /upload endpoint, parsing CSV files containing Docker container and startup time data.
// It validates the file format, ensures minimum data requirements, and stores the data for later computation.
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

// HandleCompute handles the POST /compute endpoint, computing regression models and correlation statistics from the uploaded data.
// It samples data, builds a correlation table, and computes linear, piecewise, and exponential regression models.
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
		http.Error(w, "Incorrect data provided", http.StatusBadRequest)
		return
	}

	linearRegression := table.ComputeLinearRegression()
	exponentialRegression, err := table.ComputeExponentialRegression()
	if err != nil {
		http.Error(w, "Cannot compute exponential regression model with the provided data", http.StatusBadRequest)
		return
	}

	piecewiseRegression, err := table.ComputePiecewiseRegression()
	if err != nil {
		http.Error(w, "Cannot compute piecewise regression model with the provided data", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.linearRegression = linearRegression
	s.piecewiseRegression = piecewiseRegression
	s.exponentialRegression = exponentialRegression
	s.table = table
	s.mu.Unlock()

	if err := s.writeJSONResponse(w, http.StatusOK, NewComputeResponse(sample, table,
		linearRegression, piecewiseRegression, exponentialRegression)); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// HandleSignificance handles the POST /compute/significance endpoint, performing statistical tests on previously computed models.
// It calculates Fisher F-test results for each regression model and Pearson correlation significance at the specified significance level.
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
	table, linearRegression, piecewiseRegression, exponentialRegression := s.table, s.linearRegression, s.piecewiseRegression, s.exponentialRegression
	s.mu.RUnlock()

	if table == nil {
		http.Error(w, "No computation loaded, call /compute first", http.StatusBadRequest)
		return
	}

	fisherLinear := table.ComputeFisherStatistics(linearRegression.Predict, request.SignificanceLevel, 2)
	fisherPiecewise := table.ComputeFisherStatistics(piecewiseRegression.Predict, request.SignificanceLevel, 5)
	fisherExponential := table.ComputeFisherStatistics(exponentialRegression.Predict, request.SignificanceLevel, 2)
	pearson := table.ComputePearsonCorrelation(request.SignificanceLevel)

	if err := s.writeJSONResponse(w, http.StatusOK, &SignificanceResponse{
		FisherLinear:      fisherLinear,
		FisherPiecewise:   fisherPiecewise,
		FisherExponential: fisherExponential,
		Pearson:           pearson,
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
