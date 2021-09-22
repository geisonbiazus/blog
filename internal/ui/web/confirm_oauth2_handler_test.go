package web_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/internal/ui/web"
	"github.com/geisonbiazus/blog/pkg/assert"
	"github.com/geisonbiazus/blog/pkg/testhelper"
)

func TestConfirmOauth2Handler(t *testing.T) {
	baseURL := "http://blog.example.com"

	type fixture struct {
		handler *web.ConfirmOauth2Handler
		usecase *ConfirmOauth2UseCaseSpy
	}

	setup := func() fixture {
		usecase := &ConfirmOauth2UseCaseSpy{}
		templateRenderer := newTestTemplateRenderer()
		handler := web.NewConfirmOauth2Handler(usecase, templateRenderer, baseURL)

		return fixture{
			handler: handler,
			usecase: usecase,
		}
	}

	t.Run("It responds with 404 when the given state is invalid", func(t *testing.T) {
		f := setup()

		f.usecase.ReturnError = auth.ErrInvalidState

		res := doGetRequest(f.handler, "/login/github/confirm?state=state&code=code")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, context.Background(), f.usecase.ReceivedContext)
		assert.Equal(t, "state", f.usecase.ReceivedState)
		assert.Equal(t, "code", f.usecase.ReceivedCode)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Contains(t, body, "Page not found")
	})

	t.Run("It responds with 500 when an unrecognized error is returned", func(t *testing.T) {
		f := setup()

		f.usecase.ReturnError = errors.New("some error")

		res := doGetRequest(f.handler, "/login/github/confirm?state=state&code=code")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, context.Background(), f.usecase.ReceivedContext)
		assert.Equal(t, "state", f.usecase.ReceivedState)
		assert.Equal(t, "code", f.usecase.ReceivedCode)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Contains(t, body, "Internal server error")
	})

	t.Run("It responds with redirect and setds a token cookie", func(t *testing.T) {
		f := setup()

		f.usecase.ReturnToken = "token"

		res := doGetRequest(f.handler, "/login/github/confirm?state=state&code=code")

		assert.Equal(t, context.Background(), f.usecase.ReceivedContext)
		assert.Equal(t, "state", f.usecase.ReceivedState)
		assert.Equal(t, "code", f.usecase.ReceivedCode)

		assert.Equal(t, http.StatusSeeOther, res.StatusCode)
		assert.Equal(t, baseURL, res.Header.Get("Location"))
		assert.Equal(t, "_blog_session=token", res.Cookies()[0].String())
	})
}

type ConfirmOauth2UseCaseSpy struct {
	ReceivedContext context.Context
	ReceivedState   string
	ReceivedCode    string
	ReturnToken     string
	ReturnError     error
}

func (s *ConfirmOauth2UseCaseSpy) Run(ctx context.Context, state, code string) (string, error) {
	s.ReceivedContext = ctx
	s.ReceivedState = state
	s.ReceivedCode = code
	return s.ReturnToken, s.ReturnError
}
