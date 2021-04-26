package web_test

import (
	"path/filepath"

	"github.com/geisonbiazus/blog/internal/ui/web"
)

func newTestTemplateRenderer() *web.TemplateRenderer {
	templatePath := filepath.Join("..", "..", "..", "web", "template")
	templateRenderer, err := web.NewTemplateRenderer(templatePath)

	if err != nil {
		panic(err)
	}

	return templateRenderer
}
