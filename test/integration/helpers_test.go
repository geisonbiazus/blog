package integration_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/geisonbiazus/blog/internal/app"
)

func newServer() *httptest.Server {
	os.Setenv("ENV", "test")
	c := app.NewContext()

	basePath := filepath.Join("..", "..")

	c.TemplatePath = filepath.Join(basePath, "web", "template")
	c.StaticPath = filepath.Join(basePath, "web", "static")
	c.PostPath = filepath.Join(basePath, "test", "posts")
	c.GitHubClientID = "github_client_id"
	c.GitHubClientSecret = "github_client_secret"
	return httptest.NewServer(c.Router())
}

func newNoRedirectClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
