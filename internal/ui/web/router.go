package web

import (
	"net/http"
)

func NewRouter(templatePath, staticFilesPath string, usecases *UseCases, baseURL string) http.Handler {
	templateRenderer := NewTemplateRenderer(templatePath, baseURL)

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(staticFilesPath))))
	mux.Handle("/", NewListPostsHandler(usecases.ListPosts, templateRenderer))
	mux.Handle("/posts/", NewViewPostHandler(usecases.ViewPost, templateRenderer))
	mux.Handle("/feed.atom", NewFeedHandler(usecases.ListPosts, templateRenderer, baseURL))
	mux.Handle("/about", NewTemplateHandler(templateRenderer, "about.html"))
	mux.Handle("/login/github", NewRequestOAuth2Handler(usecases.RequestOAuth2, templateRenderer))
	mux.Handle("/login/github/confirm", NewConfirmOAuth2Handler(usecases.ConfirmOAuth2, templateRenderer, baseURL))

	return mux
}
