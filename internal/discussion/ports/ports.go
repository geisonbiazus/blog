package ports

import (
	"context"

	"github.com/geisonbiazus/blog/internal/discussion/entities"
)

type CommentRepo interface {
	SaveAuthor(ctx context.Context, author *entities.Author) error
	GetAuthorByID(ctx context.Context, id string) (*entities.Author, error)
	GetAuthorByUserID(ctx context.Context, userID string) (*entities.Author, error)
	GetCommentsAndRepliesRecursively(ctx context.Context, subjectID string) ([]*entities.Comment, error)
}
