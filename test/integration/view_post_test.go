package integration_test

import (
	"net/http"
	"testing"

	"github.com/geisonbiazus/blog/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestViewPostIntegration(t *testing.T) {
	t.Run("Given a valid post path it responds with the post HTML", func(t *testing.T) {
		server := newServer()
		defer server.Close()

		res, _ := http.Get(server.URL + "/posts/test-post")

		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, body, "Test Post")
		assert.Contains(t, body, "Geison Biazus")
		assert.Contains(t, body, "April 5, 2021")
		assert.Contains(t, body, "Content")
	})
}
