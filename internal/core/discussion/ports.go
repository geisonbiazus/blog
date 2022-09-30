package discussion

import "context"

type CommentLoader interface {
	GetCommentsBySubjectID(ctx context.Context, subjectID string) ([]*Comment, error)
}

type CommentRepo interface {
	CommentLoader
}
