package web_test

import (
	"net/http"
	"testing"

	"github.com/geisonbiazus/blog/internal/ui/web"
	"github.com/geisonbiazus/blog/pkg/assert"
	"github.com/geisonbiazus/blog/pkg/testhelper"
)

func TestTemplateHandler(t *testing.T) {
	t.Run("It renders the template with the given name", func(t *testing.T) {
		templateRenderer := newTestTemplateRenderer()
		handler := web.NewTemplateHandler(templateRenderer, "about.html")

		res := doGetRequest(handler, "/")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, body, "Geison Biazus")
		assert.Contains(t, body, "Hello")
	})
}
