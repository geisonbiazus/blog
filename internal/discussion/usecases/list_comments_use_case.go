package usecases

import (
	"context"

	"github.com/geisonbiazus/blog/internal/discussion/entities"
	"github.com/geisonbiazus/blog/internal/discussion/ports"
)

type ListCommentsUseCase struct {
	commentRepo ports.CommentRepo
}

func NewListCommentsUseCase(commentRepo ports.CommentRepo) *ListCommentsUseCase {
	return &ListCommentsUseCase{commentRepo}
}

func (u *ListCommentsUseCase) Run(ctx context.Context, subjectID string) ([]*entities.Comment, error) {
	return u.commentRepo.GetCommentsAndRepliesRecursively(ctx, subjectID)
}
