package web

import (
	"context"

	"github.com/geisonbiazus/blog/internal/core/blog"
)

type UseCases struct {
	ViewPost      ViewPostUseCase
	ListPosts     ListPostUseCase
	RequestOauth2 RequestOauth2UseCase
	ConfirmOauth2 ConfirmOauth2UseCase
}

type ViewPostUseCase interface {
	Run(path string) (blog.RenderedPost, error)
}

type ListPostUseCase interface {
	Run() ([]blog.RenderedPost, error)
}

type RequestOauth2UseCase interface {
	Run() (string, error)
}

type ConfirmOauth2UseCase interface {
	Run(ctx context.Context, state, code string) (string, error)
}
