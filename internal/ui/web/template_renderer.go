package web

import (
	"html/template"
	"io"
	"path/filepath"
)

type TemplateRenderer struct {
	basePath        string
	tmpl            *template.Template
	cachedTemplates map[string]*template.Template
}

func NewTemplateRenderer(basePath string) (*TemplateRenderer, error) {
	tmpl, err := template.ParseFiles(filepath.Join(basePath, "layout.html"))
	return &TemplateRenderer{
		tmpl:            tmpl,
		basePath:        basePath,
		cachedTemplates: map[string]*template.Template{},
	}, err
}

func (r *TemplateRenderer) Render(writer io.Writer, templateName string, data interface{}) {
	tmpl := r.resolveTemplate(templateName)
	tmpl.Execute(writer, data)
}

func (r *TemplateRenderer) resolveTemplate(name string) *template.Template {
	tmpl, ok := r.cachedTemplates[name]

	if !ok {
		tmpl = r.parseTemplate(name)
		r.cachedTemplates[name] = tmpl
	}

	return tmpl
}

func (r *TemplateRenderer) parseTemplate(name string) *template.Template {
	tmpl, err := template.ParseFiles(
		filepath.Join(r.basePath, "layout.html"),
		filepath.Join(r.basePath, name),
	)

	if err != nil {
		panic(err)
	}

	return tmpl
}
