package web

import "net/http"

type ListPostHandler struct {
	usecase  ListPostUseCase
	template *TemplateRenderer
}

func NewListPostHandler(usecase ListPostUseCase, templateRenderer *TemplateRenderer) *ListPostHandler {
	return &ListPostHandler{
		usecase:  usecase,
		template: templateRenderer,
	}
}

func (h *ListPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
