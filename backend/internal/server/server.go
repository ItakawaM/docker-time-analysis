package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/ItakawaM/docker-time-analysis/internal/compute"
	"github.com/ItakawaM/docker-time-analysis/internal/parse"
)

const (
	// JSON specifies the MIME type for JSON content.
	JSON string = "application/json"
	// FormData specifies the MIME type for multipart form data.
	FormData string = "multipart/form-data"
)

// Server handles HTTP requests for data upload, computation, and significance testing.
// It maintains thread-safe state for loaded data and computed regression models.
type Server struct {
	mu  sync.RWMutex
	mux *http.ServeMux
	// set after /upload
	data []*parse.DockerEntry

	// set after /compute
	table                 *compute.CorrelationTable
	linearRegression      compute.LinearRegression
	piecewiseRegression   compute.PiecewiseRegression
	exponentialRegression compute.ExponentialRegression
}

// NewServer creates and returns a new Server instance with an initialized HTTP multiplexer.
func NewServer() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

// InitAndServe initializes HTTP routes and starts listening for requests on the specified port.
// It registers handlers for /upload, /compute, /compute/significance endpoints, and serves static files.
func (s *Server) InitAndServe(port int) {
	s.mux.HandleFunc("POST /upload", s.HandleUpload)
	s.mux.HandleFunc("POST /compute", s.HandleCompute)
	s.mux.HandleFunc("POST /compute/significance", s.HandleSignificance)

	// Serve static files from the dist directory
	distPath := "/app/dist"
	if _, err := os.Stat(distPath); err == nil {
		fileServer := http.FileServer(http.Dir(distPath))
		// Serve static assets with the /assets path
		s.mux.Handle("GET /assets/", fileServer)
		// Serve index.html for root and all other paths (SPA fallback)
		s.mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			// Try to serve the file if it exists
			filePath := filepath.Join(distPath, r.URL.Path)
			if _, err := os.Stat(filePath); err == nil && r.URL.Path != "/" {
				fileServer.ServeHTTP(w, r)
				return
			}
			// Otherwise, serve index.html for SPA routing
			http.ServeFile(w, r, filepath.Join(distPath, "index.html"))
		})
	} else {
		log.Printf("Warning: dist directory not found at %s\n", distPath)
	}

	log.Printf("Listening on :%d\n", port)

	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, s.mux))
}
