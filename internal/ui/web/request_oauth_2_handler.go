package web

import "net/http"

type RequestOauth2Handler struct {
	usecase RequestOauth2UseCase
}

func NewRequestOauth2Handler(usecase RequestOauth2UseCase) *RequestOauth2Handler {
	return &RequestOauth2Handler{usecase: usecase}
}

func (h *RequestOauth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	redirectURL := h.usecase.Run()
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
