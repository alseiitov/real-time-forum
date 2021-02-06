package server

import (
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
			Addr:    ":" + conf.GetBackendPort(),
			Handler: router,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
