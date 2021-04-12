package web

import (
	"net/http"
)

func NewRouter(templateRenderer *TemplateRenderer, usecases *UseCases) http.Handler {
	mux := http.NewServeMux()
	viewPostHandler := NewViewPostHandler(usecases.ViewPost, templateRenderer)

	mux.Handle("/", viewPostHandler)

	return mux
}
