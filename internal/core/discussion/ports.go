package discussion

import "context"

type CommentLoader interface {
	GetCommentsBySubjectID(ctx context.Context, subjectID string) ([]*Comment, error)
	GetAuthorByID(ctx context.Context, id string) (Author, error)
}

type CommentRepo interface {
	CommentLoader
}
