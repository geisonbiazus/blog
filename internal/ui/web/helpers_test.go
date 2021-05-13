package web_test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"

	"github.com/geisonbiazus/blog/internal/ui/web"
)

func newTestTemplateRenderer() *web.TemplateRenderer {
	templatePath := filepath.Join("..", "..", "..", "web", "template")
	baseURL := "http://example.com"
	templateRenderer := web.NewTemplateRenderer(templatePath, baseURL)

	return templateRenderer
}

func doGetRequest(handler http.Handler, path string) *http.Response {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rw := httptest.NewRecorder()

	handler.ServeHTTP(rw, req)

	return rw.Result()
}
