package web

import (
	"fmt"
	"net/http"
)

type Server struct {
	Port   int
	Router http.Handler
}

func NewServer(port int, router http.Handler) *Server {
	return &Server{
		Port:   port,
		Router: router,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%v", s.Port)
	return http.ListenAndServe(addr, s.Router)
}
