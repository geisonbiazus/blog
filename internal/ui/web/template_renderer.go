package web

import (
	"html/template"
	"io"
	"path/filepath"
)

type TemplateRenderer struct {
	basePath string
	tmpl     *template.Template
}

func NewTemplateRenderer(basePath string) (*TemplateRenderer, error) {
	tmpl, err := template.ParseFiles(filepath.Join(basePath, "layout.html"))
	return &TemplateRenderer{tmpl: tmpl, basePath: basePath}, err
}

func (r *TemplateRenderer) Render(writer io.Writer, templateName string, data interface{}) {
	tmpl, err := template.ParseFiles(filepath.Join(r.basePath, "layout.html"), filepath.Join(r.basePath, templateName))

	if err != nil {
		panic(err)
	}

	tmpl.Execute(writer, data)
}
