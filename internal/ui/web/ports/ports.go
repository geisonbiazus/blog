package ports

import (
	"context"

	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/internal/core/discussion"
)

type UseCases struct {
	ViewPost      ViewPostUseCase
	ListPosts     ListPostUseCase
	RequestOAuth2 RequestOAuth2UseCase
	ConfirmOAuth2 ConfirmOAuth2UseCase
	ListComments  ListCommentsUseCase
}

type ViewPostUseCase interface {
	Run(path string) (blog.RenderedPost, error)
}

type ListPostUseCase interface {
	Run() ([]blog.RenderedPost, error)
}

type RequestOAuth2UseCase interface {
	Run() (string, error)
}

type ConfirmOAuth2UseCase interface {
	Run(ctx context.Context, state, code string) (string, error)
}

type ListCommentsUseCase interface {
	Run(ctx context.Context, subjectID string) ([]*discussion.Comment, error)
}
