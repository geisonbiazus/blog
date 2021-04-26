package web

import "net/http"

type TemplateHandler struct {
	template     *TemplateRenderer
	templateName string
}

func NewTemplateHandler(renderer *TemplateRenderer, templateName string) *TemplateHandler {
	return &TemplateHandler{template: renderer, templateName: templateName}
}

func (h *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	h.template.Render(w, h.templateName, nil)
}
