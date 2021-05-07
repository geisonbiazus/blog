package integration_test

import (
	"net/http"
	"testing"

	"github.com/geisonbiazus/blog/pkg/assert"
	"github.com/geisonbiazus/blog/pkg/testhelper"
)

func TestFeedIntegration(t *testing.T) {
	t.Run("Returns the feed of the published posts", func(t *testing.T) {
		server := newServer()
		defer server.Close()

		res, _ := http.Get(server.URL + "/feed.atom")

		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/atom+xml", res.Header.Get("Content-Type"))
		assert.Contains(t, body, "<title>Geison Biazus</title>")
		assert.Contains(t, body, "<title>Test Post</title>")
		assert.Contains(t, body, "<updated>2021-04-05T18:47:00Z</updated>")
		assert.Contains(t, body, "Content")
	})
}
