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
		assert.Matches(t, res.Header.Get("Location"), `^https://github.com/login/oauth/authorize\?client_id=github_client_id&response_type=code&state=.{36}$`)
	})
}
