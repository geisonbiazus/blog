package goldmark

import (
	"bytes"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
)

type GoldmarkRenderer struct{}

func NewGoldmarkRenderer() *GoldmarkRenderer {
	return &GoldmarkRenderer{}
}

func (r *GoldmarkRenderer) Render(content string) (string, error) {
	var buf bytes.Buffer

	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					html.TabWidth(2),
				),
			),
		),
	)

	err := markdown.Convert([]byte(content), &buf)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
