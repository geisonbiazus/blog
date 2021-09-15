package web_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/geisonbiazus/blog/internal/ui/web"
	"github.com/geisonbiazus/blog/pkg/assert"
	"github.com/geisonbiazus/blog/pkg/testhelper"
)

func TestRequestOauth2Handler(t *testing.T) {
	t.Run("It requests a new Oauth2 and requiretcs to the returned URL", func(t *testing.T) {
		templateRenderer := newTestTemplateRenderer()
		usecase := &requestOauth2UseCaseStub{}
		handler := web.NewRequestOauth2Handler(usecase, templateRenderer)

		url := "http://example.com/login"

		usecase.ReturnAuthURL = url

		res := doGetRequest(handler, "/login/github/request")

		assert.Equal(t, http.StatusSeeOther, res.StatusCode)
		assert.Equal(t, url, res.Header.Get("Location"))
	})

	t.Run("It responds with 500 if an error is returned from the use case", func(t *testing.T) {
		templateRenderer := newTestTemplateRenderer()
		usecase := &requestOauth2UseCaseStub{}
		handler := web.NewRequestOauth2Handler(usecase, templateRenderer)

		usecase.ReturnError = errors.New("error")

		res := doGetRequest(handler, "/login/github/request")

		body := testhelper.ReadResponseBody(res)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.True(t, strings.Contains(body, "Internal server error"))
	})
}

type requestOauth2UseCaseStub struct {
	ReturnAuthURL string
	ReturnError   error
}

func (u requestOauth2UseCaseStub) Run() (string, error) {
	return u.ReturnAuthURL, u.ReturnError
}
