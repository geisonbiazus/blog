package handlers

import (
	"net/http"

	"github.com/geisonbiazus/blog/internal/ui/web/lib"
	"github.com/geisonbiazus/blog/internal/ui/web/ports"
)

type RequestOAuth2Handler struct {
	usecase  ports.RequestOAuth2UseCase
	template *lib.TemplateRenderer
}

func NewRequestOAuth2Handler(usecase ports.RequestOAuth2UseCase, template *lib.TemplateRenderer) *RequestOAuth2Handler {
	return &RequestOAuth2Handler{usecase: usecase, template: template}
}

func (h *RequestOAuth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	redirectURL, err := h.usecase.Run()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.template.Render(w, "500.html", nil)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
