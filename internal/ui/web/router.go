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

	mux.Handle("/posts", NewListPostsHandler(usecases.ListPosts, templateRenderer))
	mux.Handle("/", NewViewPostHandler(usecases.ViewPost, templateRenderer))

	return mux, nil
}
