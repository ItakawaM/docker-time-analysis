package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ItakawaM/docker-time-analysis/internal/parse"
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
