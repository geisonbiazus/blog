package handlers

import (
	"errors"
	"net/http"

	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/internal/ui/web/lib"
	"github.com/geisonbiazus/blog/internal/ui/web/ports"
)

type ConfirmOAuth2Handler struct {
	usecase  ports.ConfirmOAuth2UseCase
	template *lib.TemplateRenderer
	baseURL  string
}

func NewConfirmOAuth2Handler(
	usecase ports.ConfirmOAuth2UseCase,
	templateRenderer *lib.TemplateRenderer,
	baseURL string,
) *ConfirmOAuth2Handler {
	return &ConfirmOAuth2Handler{
		usecase:  usecase,
		template: templateRenderer,
		baseURL:  baseURL,
	}
}

func (h *ConfirmOAuth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	token, err := h.usecase.Run(r.Context(), state, code)

	if err != nil {
		h.respondWithError(w, err)
		return
	}

	http.SetCookie(w, h.newSessionCookie(token))
	http.Redirect(w, r, h.baseURL, http.StatusSeeOther)
}

func (h *ConfirmOAuth2Handler) respondWithError(w http.ResponseWriter, err error) {
	if errors.Is(err, auth.ErrInvalidState) {
		w.WriteHeader(http.StatusNotFound)
		h.template.Render(w, "404.html", nil)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		h.template.Render(w, "500.html", nil)
	}
}

func (h *ConfirmOAuth2Handler) newSessionCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:  "_blog_session",
		Value: token,
		Path:  "/",
	}
}
