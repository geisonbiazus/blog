package web

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Port   int
	Router http.Handler
	Logger *log.Logger
}

func NewServer(port int, router http.Handler, logger *log.Logger) *Server {
	return &Server{
		Port:   port,
		Router: router,
		Logger: logger,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%v", s.Port)
	s.Logger.Printf("Starting server on port %v", s.Port)
	return http.ListenAndServe(addr, s.Router)
}
