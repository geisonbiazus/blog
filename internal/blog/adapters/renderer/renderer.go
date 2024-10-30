package renderer

import "github.com/geisonbiazus/blog/internal/blog/adapters/renderer/goldmark"

func NewGoldmarkRenderer() *goldmark.Renderer {
	return goldmark.NewRenderer()
}
