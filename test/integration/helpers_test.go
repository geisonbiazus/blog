package integration_test

import (
	"net/http/httptest"
	"path/filepath"

	"github.com/geisonbiazus/blog/internal/app"
)

func newServer() *httptest.Server {
	c := app.NewContext()

	basePath := filepath.Join("..", "..")

	c.TemplatePath = filepath.Join(basePath, "web", "template")
	c.StaticPath = filepath.Join(basePath, "web", "static")
	c.PostPath = filepath.Join(basePath, "test", "posts")
	return httptest.NewServer(c.Router())
}
