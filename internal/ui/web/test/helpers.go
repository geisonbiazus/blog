package test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"

	"github.com/geisonbiazus/blog/internal/ui/web/lib"
)

func NewTestTemplateRenderer() *lib.TemplateRenderer {
	templatePath := filepath.Join("..", "..", "..", "..", "web", "template")
	baseURL := "http://example.com"
	templateRenderer := lib.NewTemplateRenderer(templatePath, baseURL)

	return templateRenderer
}

func DoGetRequest(handler http.Handler, path string) *http.Response {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rw := httptest.NewRecorder()

	handler.ServeHTTP(rw, req)

	return rw.Result()
}
