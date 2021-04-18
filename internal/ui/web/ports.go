package web

import "github.com/geisonbiazus/blog/internal/core/posts"

type UseCases struct {
	ViewPost ViewPostUseCase
}

type ViewPostUseCase interface {
	Run(path string) (posts.RenderedPost, error)
}

type ListPostUseCase interface {
	Run() ([]posts.Post, error)
}
