package server

import (
	"log"
	"net/http"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(conf *config.Conf, router *gorouter.Router) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + conf.BackendPort(),
			Handler: router,
		},
	}
}

func (s *Server) Run() error {
	log.Printf("API server is starting at %v", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
