package web

import (
	"fmt"
	"html/template"
	"io"
)

type TemplateRenderer struct {
	tmpl *template.Template
}

func NewTemplateRenderer(basePath string) (*TemplateRenderer, error) {
	path := fmt.Sprintf("%s/*", basePath)
	tmpl, err := template.ParseGlob(path)
	return &TemplateRenderer{tmpl: tmpl}, err
}

func (r *TemplateRenderer) Render(writer io.Writer, templateName string, data interface{}) {
	r.tmpl.ExecuteTemplate(writer, templateName, data)
}
