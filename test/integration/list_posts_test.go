package integration_test

import (
	"net/http"
	"testing"

	"github.com/geisonbiazus/blog/pkg/assert"
	"github.com/geisonbiazus/blog/pkg/testhelper"
)

func TestListPostsIntegration(t *testing.T) {
	t.Run("Returns a list of the published posts", func(t *testing.T) {
		server := newServer()
		defer server.Close()

		res, _ := http.Get(server.URL + "/posts")

		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, body, "Test Post")
		assert.Contains(t, body, "/posts/test-post")
		assert.Contains(t, body, "Geison Biazus")
		assert.Contains(t, body, "April 5, 2021")
	})
}
