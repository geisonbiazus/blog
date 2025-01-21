package ports

import (
	"context"

	"github.com/geisonbiazus/blog/internal/discussion/entities"
)

// Driving ports

type SaveAuthorInput struct {
	UserID    string
	Name      string
	AvatarURL string
}

type SaveAuthorUseCase interface {
	Run(ctx context.Context, input SaveAuthorInput) (*entities.Author, error)
}

type ListCommentsUseCase interface {
	Run(ctx context.Context, subjectID string) ([]*entities.Comment, error)
}

// Driven ports

type CommentRepo interface {
	SaveAuthor(ctx context.Context, author *entities.Author) error
	GetAuthorByID(ctx context.Context, id string) (*entities.Author, error)
	GetAuthorByUserID(ctx context.Context, userID string) (*entities.Author, error)
	GetCommentsAndRepliesRecursively(ctx context.Context, subjectID string) ([]*entities.Comment, error)
}
