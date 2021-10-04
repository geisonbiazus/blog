package web

import "net/http"

type RequestOAuth2Handler struct {
	usecase  RequestOAuth2UseCase
	template *TemplateRenderer
}

func NewRequestOAuth2Handler(usecase RequestOAuth2UseCase, template *TemplateRenderer) *RequestOAuth2Handler {
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
