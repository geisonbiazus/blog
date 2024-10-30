package goldmark

import (
	"bytes"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	htmloptions "github.com/yuin/goldmark/renderer/html"
)

type Renderer struct{}

func NewRenderer() *Renderer {
	return &Renderer{}
}

func (r *Renderer) Render(content string) (string, error) {
	var buf bytes.Buffer

	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					html.TabWidth(2),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			htmloptions.WithUnsafe(),
		),
	)

	err := markdown.Convert([]byte(content), &buf)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
