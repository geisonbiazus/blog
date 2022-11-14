package discussion

import "context"

type CommentRepo interface {
	GetCommentsAndRepliesRecursively(ctx context.Context, subjectID string) ([]*Comment, error)
}
