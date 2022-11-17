package handlers_test

import (
	"net/http"
	"testing"

	"github.com/geisonbiazus/blog/internal/ui/web/handlers"
	"github.com/geisonbiazus/blog/internal/ui/web/test"
	"github.com/geisonbiazus/blog/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestTemplateHandler(t *testing.T) {
	t.Run("It renders the template with the given name", func(t *testing.T) {
		templateRenderer := test.NewTestTemplateRenderer()
		handler := handlers.NewTemplateHandler(templateRenderer, "about.html")

		res := test.DoGetRequest(handler, "/")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, body, "Geison Biazus")
		assert.Contains(t, body, "Hello")
	})
}
