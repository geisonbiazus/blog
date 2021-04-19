package integration_test

import (
	"net/http/httptest"

	"github.com/geisonbiazus/blog/internal/app"
)

func newServer() *httptest.Server {
	c := app.NewContext()
	c.TemplatePath = "../../web/template"
	c.PostPath = "../posts"
	return httptest.NewServer(c.Router())
}
