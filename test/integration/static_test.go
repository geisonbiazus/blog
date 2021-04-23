package integration_test

import (
	"net/http"
	"testing"

	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestStatic(t *testing.T) {
	t.Run("It serves static files", func(t *testing.T) {
		server := newServer()
		defer server.Close()

		res, _ := http.Get(server.URL + "/static/styles.css")

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
