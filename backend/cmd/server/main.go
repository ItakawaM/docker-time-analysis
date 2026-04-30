package main

import (
	"github.com/ItakawaM/docker-time-analysis/internal/server"
)

func main() {
	s := server.NewServer()
	s.InitAndServe(8080)
}
