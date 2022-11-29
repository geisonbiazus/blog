package discussion

import (
	"context"
	"fmt"
)

type SaveAuthorInput struct {
	ID        string
	Name      string
	AvatarURL string
}

type SaveAuthorUseCase struct {
	commentRepo CommentRepo
}

func NewSaveAuthorUseCase(commentRepo CommentRepo) *SaveAuthorUseCase {
	return &SaveAuthorUseCase{commentRepo: commentRepo}
}

func (u *SaveAuthorUseCase) Run(ctx context.Context, input SaveAuthorInput) (*Author, error) {
	author, err := u.findOrInitializeAuthor(ctx, input.ID)
	if err != nil {
		return &Author{}, err
	}

	author.Name = input.Name
	author.AvatarURL = input.AvatarURL

	err = u.commentRepo.SaveAuthor(ctx, author)
	if err != nil {
		return &Author{}, fmt.Errorf("error on SaveAuthorUseCase.Run when saving author: %w", err)
	}

	return author, nil
}

func (u *SaveAuthorUseCase) findOrInitializeAuthor(ctx context.Context, id string) (*Author, error) {
	author, err := u.commentRepo.GetAuthorByID(ctx, id)
	if err != nil {
		return &Author{}, fmt.Errorf("error on SaveAuthorUseCase.Run when finding author: %w", err)
	}

	if author == nil {
		author = &Author{ID: id}
	}

	return author, nil
}
