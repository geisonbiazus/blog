package web

import (
	"fmt"
	"net/http"
)

type Server struct {
	Port            int
	TemplatePath    string
	ViewPostUseCase ViewPostUseCase
}

func NewServer(port int, templatePath string, viewPostUseCase ViewPostUseCase) *Server {
	return &Server{Port: port, TemplatePath: templatePath, ViewPostUseCase: viewPostUseCase}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%v", s.Port)
	templateRenderer, err := NewTemplateRenderer(s.TemplatePath)

	if err != nil {
		return err
	}

	router := NewRouter(templateRenderer, s.ViewPostUseCase)

	return http.ListenAndServe(addr, router)
}
