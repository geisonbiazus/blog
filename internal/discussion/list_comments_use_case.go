package discussion

import "context"

type ListCommentsUseCase struct {
	commentRepo CommentRepo
}

func NewListCommentsUseCase(commentRepo CommentRepo) *ListCommentsUseCase {
	return &ListCommentsUseCase{commentRepo}
}

func (u *ListCommentsUseCase) Run(ctx context.Context, subjectID string) ([]*Comment, error) {
	return u.commentRepo.GetCommentsAndRepliesRecursively(ctx, subjectID)
}
