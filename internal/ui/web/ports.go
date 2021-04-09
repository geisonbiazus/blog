package web

import "github.com/geisonbiazus/blog/internal/core/posts"

type ViewPostUseCase interface {
	Run(path string) (posts.RenderedPost, error)
}
