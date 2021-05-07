package web

import "github.com/geisonbiazus/blog/internal/core/blog"

type UseCases struct {
	ViewPost  ViewPostUseCase
	ListPosts ListPostUseCase
}

type ViewPostUseCase interface {
	Run(path string) (blog.RenderedPost, error)
}

type ListPostUseCase interface {
	Run() ([]blog.RenderedPost, error)
}
