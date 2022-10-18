package discussion

import (
	"context"
	"time"
)

type CommentParams struct {
	ID        string
	SubjectID string
	AuthorID  string
	Markdown  string
	HTML      string
	CreatedAt time.Time
}

type Comment struct {
	ID        string
	SubjectID string
	AuthorID  string
	Markdown  string
	HTML      string
	CreatedAt time.Time

	author  Author
	replies []*Comment

	loader CommentLoader
}

func NewComment(params CommentParams, commentLoader CommentLoader) *Comment {
	return &Comment{
		ID:        params.ID,
		SubjectID: params.SubjectID,
		AuthorID:  params.AuthorID,
		Markdown:  params.Markdown,
		HTML:      params.HTML,
		CreatedAt: params.CreatedAt,
		loader:    commentLoader,
	}
}

func (c *Comment) Author(ctx context.Context) (Author, error) {
	if (c.author == Author{}) {
		author, err := c.loader.GetAuthorByID(ctx, c.AuthorID)
		if err != nil {
			return Author{}, err
		}

		c.author = author
	}

	return c.author, nil
}

func (c *Comment) SetAuthor(author Author) {
	c.author = author
}

func (c *Comment) Replies(ctx context.Context) ([]*Comment, error) {
	if c.replies == nil {
		replies, err := c.loader.GetCommentsBySubjectID(ctx, c.ID)
		if err != nil {
			return []*Comment{}, err
		}

		c.replies = replies
	}

	return c.replies, nil
}

func (c *Comment) SetReplies(replies []*Comment) {
	c.replies = replies
}
