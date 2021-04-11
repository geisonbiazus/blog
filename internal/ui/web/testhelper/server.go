package testhelper

import (
	"net/http/httptest"

	"github.com/geisonbiazus/blog/internal/ui/web"
)

func NewTestServer() (*httptest.Server, *web.UseCases) {
	templateRenderer, err := web.NewTemplateRenderer("../../../web/template")

	if err != nil {
		panic(err)
	}

	usecases := &web.UseCases{
		ViewPost: &ViewPostUseCaseSpy{},
	}
	router := web.NewRouter(templateRenderer, usecases)
	server := httptest.NewServer(router)

	return server, usecases
}
