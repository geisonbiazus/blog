package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/geisonbiazus/blog/internal/app"
	"github.com/geisonbiazus/blog/pkg/assert"
	"github.com/geisonbiazus/blog/pkg/testhelper"
)

func TestViewPostIntegration(t *testing.T) {
	setup := func() *httptest.Server {
		c := app.NewContext()
		c.TemplatePath = "../../web/template"
		c.PostPath = "../posts"
		server := httptest.NewServer(c.Router())

		return server
	}

	t.Run("Given a valid post path it responds with the post HTML", func(t *testing.T) {
		server := setup()
		defer server.Close()

		res, _ := http.Get(server.URL + "/test-post")

		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.True(t, strings.Contains(body, "Test Post"))
		assert.True(t, strings.Contains(body, "Geison Biazus"))
		assert.True(t, strings.Contains(body, "05 Apr 21"))
		assert.True(t, strings.Contains(body, "Content"))
	})
}
