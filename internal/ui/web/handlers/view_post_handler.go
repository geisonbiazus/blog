package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/geisonbiazus/blog/internal/ui/web/lib"
	"github.com/geisonbiazus/blog/internal/ui/web/ports"
)

type ViewPostHandler struct {
	viewPostUseCase     ports.ViewPostUseCase
	listCommentsUseCase ports.ListCommentsUseCase
	template            *lib.TemplateRenderer
}

func NewViewPostHandler(
	viewPostUseCase ports.ViewPostUseCase,
	listCommentsUseCase ports.ListCommentsUseCase,
	templateRenderer *lib.TemplateRenderer,
) *ViewPostHandler {
	return &ViewPostHandler{
		viewPostUseCase:     viewPostUseCase,
		listCommentsUseCase: listCommentsUseCase,
		template:            templateRenderer,
	}
}

func (h *ViewPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := path.Base(r.URL.Path)
	renderedPost, err := h.viewPostUseCase.Run(path)
	comments, _ := h.listCommentsUseCase.Run(r.Context(), path)

	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		h.template.Render(w, "view_post.html", h.toViewModel(renderedPost, comments))
	case blog.ErrPostNotFound:
		w.WriteHeader(http.StatusNotFound)
		h.template.Render(w, "404.html", nil)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		h.template.Render(w, "500.html", nil)
	}
}

func (h *ViewPostHandler) toViewModel(p blog.RenderedPost, comments []*discussion.Comment) postViewModel {
	return postViewModel{
		Title:       p.Post.Title,
		Author:      p.Post.Author,
		Description: p.Post.Description,
		ImagePath:   p.Post.ImagePath,
		Path:        fmt.Sprintf("/posts/%s", p.Post.Path),
		Date:        p.Post.Time.Format(lib.DateFormat),
		Content:     template.HTML(p.HTML),
		Comments:    h.toCommentsViewModel(comments),
	}
}

func (h *ViewPostHandler) toCommentsViewModel(comments []*discussion.Comment) []commentViewModel {
	result := []commentViewModel{}

	for _, comment := range comments {
		viewModel := commentViewModel{
			AuthorAvatarURL: comment.Author.AvatarURL,
			AuthorName:      comment.Author.Name,
			Date:            comment.CreatedAt.Format(lib.DateFormat),
			Content:         template.HTML(comment.HTML),
		}

		if comment.Replies != nil {
			viewModel.Replies = h.toCommentsViewModel(comment.Replies)
		}

		result = append(result, viewModel)
	}

	return result
}

type postViewModel struct {
	Title       string
	Author      string
	Date        string
	Description string
	ImagePath   string
	Path        string
	Content     template.HTML
	Comments    []commentViewModel
}

type commentViewModel struct {
	AuthorAvatarURL string
	AuthorName      string
	Date            string
	Content         template.HTML
	Replies         []commentViewModel
}
