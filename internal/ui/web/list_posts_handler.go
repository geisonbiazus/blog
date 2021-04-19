package web

import (
	"net/http"

	"github.com/geisonbiazus/blog/internal/core/posts"
)

type ListPostsHandler struct {
	usecase  ListPostUseCase
	template *TemplateRenderer
}

func NewListPostsHandler(usecase ListPostUseCase, templateRenderer *TemplateRenderer) *ListPostsHandler {
	return &ListPostsHandler{
		usecase:  usecase,
		template: templateRenderer,
	}
}

func (h *ListPostsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postList, err := h.usecase.Run()

	if err == nil {
		models := h.toViewModelList(postList)
		w.WriteHeader(http.StatusOK)
		h.template.Render(w, "list_posts.html", models)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		h.template.Render(w, "500.html", nil)
	}
}

func (h *ListPostsHandler) toViewModelList(postList []posts.Post) []postListViewModel {
	models := []postListViewModel{}

	for _, post := range postList {
		models = append(models, h.toViewModel(post))
	}

	return models
}

func (h *ListPostsHandler) toViewModel(post posts.Post) postListViewModel {
	return postListViewModel{
		Title:  post.Title,
		Author: post.Author,
		Date:   post.Time.Format("02 Jan 06"),
		Path:   post.Path,
	}
}

type postListViewModel struct {
	Title  string
	Path   string
	Author string
	Date   string
}
