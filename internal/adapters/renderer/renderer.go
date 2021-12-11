package renderer

import "github.com/geisonbiazus/blog/internal/adapters/renderer/goldmark"

func NewGoldmarkRenderer() *goldmark.Renderer {
	return goldmark.NewRenderer()
}
