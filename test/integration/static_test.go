package integration_test

import (
	"net/http"
	"testing"

	"github.com/geisonbiazus/blog/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestStaticIntegration(t *testing.T) {
	t.Run("It serves static files", func(t *testing.T) {
		server := newServer()
		defer server.Close()

		res, _ := http.Get(server.URL + "/static/styles.css")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, body, ".blog-container")
	})
}
