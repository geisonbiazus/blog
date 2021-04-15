package web

import (
	"net/http"
)

func NewRouter(templatePath string, usecases *UseCases) (http.Handler, error) {
	templateRenderer, err := NewTemplateRenderer(templatePath)

	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	viewPostHandler := NewViewPostHandler(usecases.ViewPost, templateRenderer)

	mux.Handle("/", viewPostHandler)

	return mux, nil
}
