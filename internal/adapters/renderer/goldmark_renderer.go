package renderer

import (
	"bytes"

	"github.com/yuin/goldmark"
)

type GoldmarkRenderer struct{}

func NewGoldmarkRenderer() *GoldmarkRenderer {
	return &GoldmarkRenderer{}
}

func (r *GoldmarkRenderer) Render(content string) (string, error) {
	var buf bytes.Buffer
	err := goldmark.Convert([]byte(content), &buf)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
