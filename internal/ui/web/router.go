package web

import (
	"net/http"
)

func NewRouter(templateRenderer *TemplateRenderer, viewPostUseCase ViewPostUseCase) http.Handler {
	mux := http.NewServeMux()
	viewPostHandler := NewViewPostHandler(viewPostUseCase, templateRenderer)

	mux.Handle("/", viewPostHandler)

	return mux
}
