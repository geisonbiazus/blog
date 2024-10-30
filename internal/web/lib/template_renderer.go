package lib

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
)

type TemplateRenderer struct {
	basePath        string
	baseURL         string
	cachedTemplates map[string]*template.Template
}

func NewTemplateRenderer(basePath, baseURL string) *TemplateRenderer {
	return &TemplateRenderer{
		basePath:        basePath,
		baseURL:         baseURL,
		cachedTemplates: map[string]*template.Template{},
	}
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

	return tmpl.Lookup("layout.html")
}

func (r *TemplateRenderer) parseTemplate(name string) *template.Template {
	tmpl, err := template.New("template").Funcs(r.templateFuncs()).ParseFiles(
		filepath.Join(r.basePath, "layout.html"),
		filepath.Join(r.basePath, name),
	)

	if err != nil {
		panic(err)
	}

	return tmpl
}

func (r *TemplateRenderer) templateFuncs() template.FuncMap {
	return template.FuncMap{
		"urlFor": r.urlFor,
	}
}

func (r *TemplateRenderer) urlFor(path string) string {
	return fmt.Sprintf("%s%s", r.baseURL, path)
}
