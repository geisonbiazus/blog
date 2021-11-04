package integration_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/geisonbiazus/blog/pkg/assert"
)

var locationPattern = `^http://localhost:3000/login/github/confirm\?state=(.{36})&code=(.{36})$`
var locationRegex = regexp.MustCompile(locationPattern)

func TestRequestOAuth2Integration(t *testing.T) {
	t.Run("It redirects to github", func(t *testing.T) {
		server := newServer()
		client := newNoRedirectClient()

		defer server.Close()

		res, _ := client.Get(server.URL + "/login/github")

		assert.Equal(t, http.StatusSeeOther, res.StatusCode)
		assert.Matches(t, locationPattern, res.Header.Get("Location"))

		matches := locationRegex.FindStringSubmatch(res.Header.Get("Location"))
		state := matches[1]
		code := matches[2]

		res, _ = client.Get(fmt.Sprintf("%s/login/github/confirm?state=%s&code=%s", server.URL, state, code))

		assert.Equal(t, http.StatusSeeOther, res.StatusCode)
		assert.Matches(t, "_blog_session=.+; Path=/", res.Cookies()[0].String())
	})
}
