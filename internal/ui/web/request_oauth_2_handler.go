package web

import "net/http"

type RequestOauth2Handler struct {
	usecase  RequestOauth2UseCase
	template *TemplateRenderer
}

func NewRequestOauth2Handler(usecase RequestOauth2UseCase, template *TemplateRenderer) *RequestOauth2Handler {
	return &RequestOauth2Handler{usecase: usecase, template: template}
}

func (h *RequestOauth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	redirectURL, err := h.usecase.Run()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.template.Render(w, "500.html", nil)
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
