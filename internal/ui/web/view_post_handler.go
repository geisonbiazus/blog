package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/geisonbiazus/blog/internal/core/blog"
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
	path := path.Base(req.URL.Path)
	renderedPost, err := h.usecase.Run(path)

	switch err {
	case nil:
		res.WriteHeader(http.StatusOK)
		h.template.Render(res, "view_post.html", h.toViewModel(renderedPost))
	case blog.ErrPostNotFound:
		res.WriteHeader(http.StatusNotFound)
		h.template.Render(res, "404.html", nil)
	default:
		res.WriteHeader(http.StatusInternalServerError)
		h.template.Render(res, "500.html", nil)
	}
}

func (h *ViewPostHandler) toViewModel(p blog.RenderedPost) postViewModel {
	return postViewModel{
		Title:       p.Post.Title,
		Author:      p.Post.Author,
		Description: p.Post.Description,
		ImagePath:   p.Post.ImagePath,
		Path:        fmt.Sprintf("/posts/%s", p.Post.Path),
		Date:        p.Post.Time.Format(DateFormat),
		Content:     template.HTML(p.HTML),
	}
}

type postViewModel struct {
	Title       string
	Author      string
	Date        string
	Description string
	ImagePath   string
	Path        string
	Content     template.HTML
}
