package discussion

import "context"

type CommentRepo interface {
	SaveAuthor(ctx context.Context, author *Author) error
	GetCommentsAndRepliesRecursively(ctx context.Context, subjectID string) ([]*Comment, error)
}
