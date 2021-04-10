package web

import "net/http"

func NewRouter(viewPostUseCase ViewPostUseCase) http.Handler {
	mux := http.NewServeMux()

	templateRenderer, _ := NewTemplateRenderer()
	viewPostHandler := NewViewPostHandler(viewPostUseCase, templateRenderer)

	mux.Handle("/", viewPostHandler)

	return mux
}
