package server

import (
	"fmt"
	"log"
	"net/http"
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
// It registers handlers for /upload, /compute, and /compute/significance endpoints.
func (s *Server) InitAndServe(port int) {
	s.mux.HandleFunc("POST /upload", s.HandleUpload)
	s.mux.HandleFunc("POST /compute", s.HandleCompute)
	s.mux.HandleFunc("POST /compute/significance", s.HandleSignificance)

	log.Printf("Listening on :%d\n", port)

	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, s.mux))
}
