package web

import (
	"net/http"
)

func NewRouter(templatePath, staticFilesPath string, usecases *UseCases) (http.Handler, error) {
	templateRenderer, err := NewTemplateRenderer(templatePath)

	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(staticFilesPath))))
	mux.Handle("/posts", NewListPostsHandler(usecases.ListPosts, templateRenderer))
	mux.Handle("/posts/", NewViewPostHandler(usecases.ViewPost, templateRenderer))
	mux.Handle("/", NewTemplateHandler(templateRenderer, "home.html"))

	return mux, nil
}
