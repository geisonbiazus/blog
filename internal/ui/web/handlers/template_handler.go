package handlers

import (
	"net/http"

	"github.com/geisonbiazus/blog/internal/ui/web/lib"
)

type TemplateHandler struct {
	template     *lib.TemplateRenderer
	templateName string
}

func NewTemplateHandler(renderer *lib.TemplateRenderer, templateName string) *TemplateHandler {
	return &TemplateHandler{template: renderer, templateName: templateName}
}

func (h *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	h.template.Render(w, h.templateName, nil)
}
