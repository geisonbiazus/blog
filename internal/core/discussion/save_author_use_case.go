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
	author := u.authorFrom(input)

	err := u.commentRepo.SaveAuthor(ctx, author)
	if err != nil {
		return &Author{}, fmt.Errorf("error on SaveAuthorUseCase.Run when saving author: %w", err)
	}

	return author, nil
}

func (u *SaveAuthorUseCase) authorFrom(input SaveAuthorInput) *Author {
	return &Author{
		ID:        input.ID,
		Name:      input.Name,
		AvatarURL: input.AvatarURL,
	}
}
