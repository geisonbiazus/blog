package website_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geisonbiazus/blog/internal/ui/website"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestViewPostHandler(t *testing.T) {
	t.Run("Given a valid post path it responds with the post HTML", func(t *testing.T) {
		server := httptest.NewServer(website.NewRouter())
		defer server.Close()

		res, _ := http.Get(server.URL + "/test-post")

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
