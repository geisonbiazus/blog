package discussion

import (
	"context"
	"fmt"

	"github.com/geisonbiazus/blog/internal/core/shared"
)

type SaveAuthorInput struct {
	UserID    string
	Name      string
	AvatarURL string
}

type SaveAuthorUseCase struct {
	commentRepo CommentRepo
	txManager   shared.TransactionManager
	idGen       shared.IDGenerator
}

func NewSaveAuthorUseCase(commentRepo CommentRepo, txManager shared.TransactionManager, idGen shared.IDGenerator) *SaveAuthorUseCase {
	return &SaveAuthorUseCase{
		commentRepo: commentRepo,
		txManager:   txManager,
		idGen:       idGen,
	}
}

func (u *SaveAuthorUseCase) Run(ctx context.Context, input SaveAuthorInput) (author *Author, err error) {
	u.txManager.Transaction(ctx, func(ctx context.Context) error {
		author, err = u.run(ctx, input)
		return err
	})
	return
}

func (u *SaveAuthorUseCase) run(ctx context.Context, input SaveAuthorInput) (*Author, error) {
	author, err := u.findOrInitializeAuthor(ctx, input.UserID)
	if err != nil {
		return &Author{}, err
	}

	u.setAuthorAttributes(author, input)

	err = u.commentRepo.SaveAuthor(ctx, author)
	if err != nil {
		return &Author{}, fmt.Errorf("error on SaveAuthorUseCase.Run when saving author: %w", err)
	}

	return author, nil
}

func (u *SaveAuthorUseCase) findOrInitializeAuthor(ctx context.Context, id string) (*Author, error) {
	author, err := u.commentRepo.GetAuthorByUserID(ctx, id)
	if err != nil {
		return &Author{}, fmt.Errorf("error on SaveAuthorUseCase.Run when finding author: %w", err)
	}

	if author == nil {
		author = &Author{ID: u.idGen.Generate()}
	}

	return author, nil
}

func (u *SaveAuthorUseCase) setAuthorAttributes(author *Author, input SaveAuthorInput) {
	author.UserID = input.UserID
	author.Name = input.Name
	author.AvatarURL = input.AvatarURL
}
