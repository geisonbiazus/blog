package discussion

import "context"

type CommentRepo interface {
	GetCommentsBySubjectID(ctx context.Context, subjectID string) ([]*Comment, error)
	GetAuthorByID(ctx context.Context, id string) (Author, error)
}
