package web

import (
	"net/http"

	"github.com/geisonbiazus/blog/internal/ui/web/handlers"
	"github.com/geisonbiazus/blog/internal/ui/web/lib"
	"github.com/geisonbiazus/blog/internal/ui/web/ports"
)

func NewRouter(templatePath, staticFilesPath string, usecases *ports.UseCases, baseURL string) http.Handler {
	templateRenderer := lib.NewTemplateRenderer(templatePath, baseURL)

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(staticFilesPath))))
	mux.Handle("/", handlers.NewListPostsHandler(usecases.ListPosts, templateRenderer))
	mux.Handle("/posts/", handlers.NewViewPostHandler(usecases.ViewPost, templateRenderer))
	mux.Handle("/feed.atom", handlers.NewFeedHandler(usecases.ListPosts, templateRenderer, baseURL))
	mux.Handle("/about", handlers.NewTemplateHandler(templateRenderer, "about.html"))
	mux.Handle("/login/github", handlers.NewRequestOAuth2Handler(usecases.RequestOAuth2, templateRenderer))
	mux.Handle("/login/github/confirm", handlers.NewConfirmOAuth2Handler(usecases.ConfirmOAuth2, templateRenderer, baseURL))

	return mux
}
