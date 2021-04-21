package web

import (
	"bytes"
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
	tmpl, _ := template.ParseGlob(r.path)
	r.tmpl = tmpl

	buf := bytes.NewBuffer([]byte{})
	r.tmpl.ExecuteTemplate(buf, templateName, data)
	r.tmpl.ExecuteTemplate(writer, "layout.html", template.HTML(buf.String()))
}
