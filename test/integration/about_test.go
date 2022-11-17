package integration_test

import (
	"net/http"
	"testing"

	"github.com/geisonbiazus/blog/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestHomeIntegration(t *testing.T) {
	t.Run("It renders home page", func(t *testing.T) {
		server := newServer()
		defer server.Close()

		res, _ := http.Get(server.URL + "/about")

		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, body, "Geison Biazus")
		assert.Contains(t, body, "Hello")
	})
}
