package web_test

import (
	"github.com/geisonbiazus/blog/internal/ui/web"
)

func newTestTemplateRenderer() *web.TemplateRenderer {
	templateRenderer, err := web.NewTemplateRenderer("../../../web/template")

	if err != nil {
		panic(err)
	}

	return templateRenderer
}
