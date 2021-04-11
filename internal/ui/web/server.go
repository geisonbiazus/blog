package web

import (
	"fmt"
	"net/http"
)

type Server struct {
	Port         int
	TemplatePath string
	UseCases     *UseCases
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%v", s.Port)
	templateRenderer, err := NewTemplateRenderer(s.TemplatePath)

	if err != nil {
		return err
	}

	router := NewRouter(templateRenderer, s.UseCases)

	return http.ListenAndServe(addr, router)
}
