package usecases

import (
	"context"
	"fmt"

	"github.com/geisonbiazus/blog/internal/discussion/entities"
	"github.com/geisonbiazus/blog/internal/discussion/ports"
	"github.com/geisonbiazus/blog/pkg/gen"
	"github.com/geisonbiazus/blog/pkg/transaction"
)

type SaveAuthorUseCase struct {
	commentRepo ports.CommentRepo
	txManager   transaction.Manager
	idGen       gen.Generator
}

func NewSaveAuthorUseCase(commentRepo ports.CommentRepo, txManager transaction.Manager, idGen gen.Generator) *SaveAuthorUseCase {
	return &SaveAuthorUseCase{
		commentRepo: commentRepo,
		txManager:   txManager,
		idGen:       idGen,
	}
}

func (u *SaveAuthorUseCase) Run(ctx context.Context, input ports.SaveAuthorInput) (author *entities.Author, err error) {
	u.txManager.Transaction(ctx, func(ctx context.Context) error {
		author, err = u.run(ctx, input)
		return err
	})
	return
}

func (u *SaveAuthorUseCase) run(ctx context.Context, input ports.SaveAuthorInput) (*entities.Author, error) {
	author, err := u.findOrInitializeAuthor(ctx, input.UserID)
	if err != nil {
		return &entities.Author{}, err
	}

	u.setAuthorAttributes(author, input)

	err = u.commentRepo.SaveAuthor(ctx, author)
	if err != nil {
		return &entities.Author{}, fmt.Errorf("error on SaveAuthorUseCase.Run when saving author: %w", err)
	}

	return author, nil
}

func (u *SaveAuthorUseCase) findOrInitializeAuthor(ctx context.Context, id string) (*entities.Author, error) {
	author, err := u.commentRepo.GetAuthorByUserID(ctx, id)
	if err != nil {
		return &entities.Author{}, fmt.Errorf("error on SaveAuthorUseCase.Run when finding author: %w", err)
	}

	if author == nil {
		author = &entities.Author{ID: u.idGen.Generate()}
	}

	return author, nil
}

func (u *SaveAuthorUseCase) setAuthorAttributes(author *entities.Author, input ports.SaveAuthorInput) {
	author.UserID = input.UserID
	author.Name = input.Name
	author.AvatarURL = input.AvatarURL
}
