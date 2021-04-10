package web

import (
	"html/template"
	"io"
)

type TemplateRenderer struct {
	tmpl *template.Template
}

func NewTemplateRenderer() (*TemplateRenderer, error) {
	tmpl, err := template.ParseGlob("../../../web/template/*")
	return &TemplateRenderer{tmpl: tmpl}, err
}

func (r *TemplateRenderer) Render(writter io.Writer, templateName string, data interface{}) {
	r.tmpl.ExecuteTemplate(writter, templateName, data)
}
