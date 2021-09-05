package web_test

import (
	"net/http"
	"testing"

	"github.com/geisonbiazus/blog/internal/ui/web"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestRequestOauth2Handler(t *testing.T) {
	t.Run("It requests a new Oauth2 and requiretcs to the returned URL", func(t *testing.T) {
		usecase := &requestOauth2UseCaseStub{}
		handler := web.NewRequestOauth2Handler(usecase)

		url := "http://example.com/login"

		usecase.ReturnAuthURL = url

		res := doGetRequest(handler, "/login/github/request")

		assert.Equal(t, http.StatusSeeOther, res.StatusCode)
		assert.Equal(t, url, res.Header.Get("Location"))
	})
}

type requestOauth2UseCaseStub struct {
	ReturnAuthURL string
}

func (u requestOauth2UseCaseStub) Run() string {
	return u.ReturnAuthURL
}
