package web

import (
	"errors"
	"net/http"

	"github.com/geisonbiazus/blog/internal/core/auth"
)

const _24HoursInSeconds = 24 * 60 * 60

type ConfirmOauth2Handler struct {
	usecase  ConfirmOauth2UseCase
	template *TemplateRenderer
	baseURL  string
}

func NewConfirmOauth2Handler(
	usecase ConfirmOauth2UseCase,
	templateRenderer *TemplateRenderer,
	baseURL string,
) *ConfirmOauth2Handler {
	return &ConfirmOauth2Handler{
		usecase:  usecase,
		template: templateRenderer,
		baseURL:  baseURL,
	}
}

func (h *ConfirmOauth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	token, err := h.usecase.Run(r.Context(), state, code)

	if err != nil {
		h.respondWithError(w, err)
	}

	http.SetCookie(w, h.newSessionCookie(token))
	http.Redirect(w, r, h.baseURL, http.StatusSeeOther)
}

func (h *ConfirmOauth2Handler) respondWithError(w http.ResponseWriter, err error) {
	if errors.Is(err, auth.ErrInvalidState) {
		w.WriteHeader(http.StatusNotFound)
		h.template.Render(w, "404.html", nil)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		h.template.Render(w, "500.html", nil)
	}
}

func (h *ConfirmOauth2Handler) newSessionCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:   "_blog_session",
		Value:  token,
		MaxAge: _24HoursInSeconds,
	}
}
