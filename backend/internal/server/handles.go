package server

import (
	"encoding/json"
	"net/http"

	"github.com/ItakawaM/docker-time-analysis/internal/parse"
)

type UploadResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}

func (s *Server) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file field", http.StatusBadRequest)
		return
	}

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
	json.NewEncoder(w).Encode(UploadResponse{
		Message: "Parsed the file",
		Count:   len(parsedData),
	})
}
