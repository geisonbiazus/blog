package discussion

import (
	"context"
)

type CommentParams struct {
	ID        string
	SubjectID string
	Markdown  string
	HTML      string
}

type Comment struct {
	ID        string
	SubjectID string
	Markdown  string
	HTML      string

	loader CommentLoader
}

func NewComment(params CommentParams, commentLoader CommentLoader) *Comment {
	return &Comment{
		ID:        params.ID,
		SubjectID: params.SubjectID,
		Markdown:  params.Markdown,
		HTML:      params.HTML,
		loader:    commentLoader,
	}
}

func (c *Comment) Replies(ctx context.Context) ([]*Comment, error) {
	return c.loader.GetCommentsBySubjectID(ctx, c.ID)
}
