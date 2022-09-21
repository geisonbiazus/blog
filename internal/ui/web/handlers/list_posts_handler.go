package handlers

import (
	"fmt"
	"net/http"

	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/internal/ui/web/lib"
	"github.com/geisonbiazus/blog/internal/ui/web/ports"
)

type ListPostsHandler struct {
	usecase  ports.ListPostUseCase
	template *lib.TemplateRenderer
}

func NewListPostsHandler(usecase ports.ListPostUseCase, templateRenderer *lib.TemplateRenderer) *ListPostsHandler {
	return &ListPostsHandler{
		usecase:  usecase,
		template: templateRenderer,
	}
}

func (h *ListPostsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	posts, err := h.usecase.Run()

	if err == nil {
		models := h.toViewModelList(posts)
		w.WriteHeader(http.StatusOK)
		h.template.Render(w, "list_posts.html", models)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		h.template.Render(w, "500.html", nil)
	}
}

func (h *ListPostsHandler) toViewModelList(posts []blog.RenderedPost) []postsViewModel {
	models := []postsViewModel{}

	for _, post := range posts {
		models = append(models, h.toViewModel(post))
	}

	return models
}

func (h *ListPostsHandler) toViewModel(post blog.RenderedPost) postsViewModel {
	return postsViewModel{
		Title:  post.Post.Title,
		Author: post.Post.Author,
		Date:   post.Post.Time.Format(lib.DateFormat),
		Path:   fmt.Sprintf("/posts/%s", post.Post.Path),
	}
}

type postsViewModel struct {
	Title  string
	Path   string
	Author string
	Date   string
}
