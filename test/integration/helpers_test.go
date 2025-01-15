package integration_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/geisonbiazus/blog/internal/app"
)

func newServer() *httptest.Server {
	basePath := filepath.Join("..", "..")

	os.Setenv("ENV", "test")
	os.Setenv("POST_PATH", filepath.Join(basePath, "test", "posts"))
	os.Setenv("TEMPLATE_PATH", filepath.Join(basePath, "web", "template"))
	os.Setenv("STATIC_PATH", filepath.Join(basePath, "web", "static"))
	os.Setenv("GITHUB_CLIENT_ID", "github_client_id")
	os.Setenv("GITHUB_CLIENT_SECRET", "github_client_secret")

	c := app.NewContext()

	return httptest.NewServer(c.Router())
}

func newNoRedirectClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
