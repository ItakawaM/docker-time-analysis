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
	JSON     string = "application/json"
	FormData string = "multipart/form-data"
)

type Server struct {
	mu  sync.RWMutex
	mux *http.ServeMux
	// set after /upload
	data []*parse.DockerEntry

	// set after /compute
	table *compute.CorrelationTable
	alpha float64
	beta  float64
}

func NewServer() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

func (s *Server) InitAndServe(port int) {
	s.mux.HandleFunc("POST /upload", s.HandleUpload)
	s.mux.HandleFunc("POST /compute", s.HandleCompute)
	s.mux.HandleFunc("POST /compute/significance", s.HandleSignificance)

	log.Printf("Listening on :%d\n", port)

	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, s.mux))
}
