package ports

import (
	"github.com/geisonbiazus/blog/internal/auth"
	"github.com/geisonbiazus/blog/internal/blog"
	"github.com/geisonbiazus/blog/internal/discussion"
)

type UseCases struct {
	ViewPost      blog.ViewPostUseCase
	ListPosts     blog.ListPostsUseCase
	RequestOAuth2 auth.RequestOAuth2UseCase
	ConfirmOAuth2 auth.ConfirmOAuth2UseCase
	ListComments  discussion.ListCommentsUseCase
}
