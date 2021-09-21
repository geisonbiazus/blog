package integration_test

import (
	"net/http"
	"testing"

	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestRequestOauth2Integration(t *testing.T) {
	t.Run("It redirects to github", func(t *testing.T) {
		server := newServer()
		client := newNoRedirectClient()

		defer server.Close()

		res, _ := client.Get(server.URL + "/login/github/request")

		assert.Equal(t, http.StatusSeeOther, res.StatusCode)
		assert.Matches(t, res.Header.Get("Location"), `^http://localhost:3000/login/github/confirm\?state=.{36}&code=.{36}$`)
	})
}
