package web

import (
	"fmt"
	"html/template"
	"io"
)

type TemplateRenderer struct {
	path string
	tmpl *template.Template
}

func NewTemplateRenderer(basePath string) (*TemplateRenderer, error) {
	path := fmt.Sprintf("%s/*", basePath)
	tmpl, err := template.ParseGlob(path)
	return &TemplateRenderer{tmpl: tmpl, path: path}, err
}

func (r *TemplateRenderer) Render(writer io.Writer, templateName string, data interface{}) {
	r.tmpl.ExecuteTemplate(writer, templateName, data)
}
