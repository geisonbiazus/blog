package web

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/geisonbiazus/blog/internal/core/posts"
)

type ViewPostHandler struct {
	usecase  ViewPostUseCase
	template *TemplateRenderer
}

func NewViewPostHandler(usecase ViewPostUseCase, templateRenderer *TemplateRenderer) *ViewPostHandler {
	return &ViewPostHandler{
		usecase:  usecase,
		template: templateRenderer,
	}
}

func (h *ViewPostHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/")
	renderedPost, _ := h.usecase.Run(path)

	res.WriteHeader(http.StatusOK)
	h.template.Render(res, "post.html", h.toViewModel(renderedPost))
}

func (h *ViewPostHandler) toViewModel(p posts.RenderedPost) postViewModel {
	return postViewModel{
		Title:   p.Title,
		Author:  p.Author,
		Date:    p.Time.Format("02 Jan 06"),
		Content: template.HTML(p.Content),
	}
}

type postViewModel struct {
	Title   string
	Author  string
	Date    string
	Content template.HTML
}
