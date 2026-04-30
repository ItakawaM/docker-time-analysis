package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/ItakawaM/docker-time-analysis/internal/parse"
)

type Server struct {
	mu   sync.RWMutex
	mux  *http.ServeMux
	data []*parse.DockerEntry
}

func NewServer() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

func (s *Server) InitAndServe(port int) {
	s.mux.HandleFunc("POST /upload", s.HandleUpload)
	s.mux.HandleFunc("POST /compute", s.HandleCompute)

	log.Printf("Listening on :%d\n", port)

	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, s.mux))
}
