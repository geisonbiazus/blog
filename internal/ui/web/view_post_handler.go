package web

import (
	"html/template"
	"io"
	"net/http"
	"strings"
)

type ViewPostHandler struct {
	usecase ViewPostUseCase
}

func NewViewPostHandler(usecase ViewPostUseCase) *ViewPostHandler {
	return &ViewPostHandler{
		usecase: usecase,
	}
}

func (h *ViewPostHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/")
	renderedPost, _ := h.usecase.Run(path)

	renderTemplate(res, "post.go.tmpl", renderedPost)
	res.WriteHeader(http.StatusOK)
}

func renderTemplate(writter io.Writer, templateName string, data interface{}) {
	fm := template.FuncMap{
		"htmlsafe": func(html string) template.HTML {
			return template.HTML(html)
		},
	}

	tmpl := template.Must(template.New("").Funcs(fm).ParseGlob("../../../web/template/*.tmpl"))

	tmpl.ExecuteTemplate(writter, templateName, data)
}
