package discussion_test

import (
	"context"
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/commentrepo/memory"
	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/stretchr/testify/suite"
)

type SaveAuthorUseCaseSuite struct {
	suite.Suite
	usecase *discussion.SaveAuthorUseCase
	repo    *memory.CommentRepo
	ctx     context.Context
}

func (s *SaveAuthorUseCaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = memory.NewCommentRepo()
	s.usecase = discussion.NewSaveAuthorUseCase(s.repo)
}

func (s *SaveAuthorUseCaseSuite) TestRun() {
	s.Run("It creates an author when it doesn't exist", func() {
		author, err := s.usecase.Run(s.ctx, s.input())

		s.Equal(s.author(), author)
		s.Nil(err)

		persistedAuthor, _ := s.repo.GetAuthorByID(s.ctx, s.input().ID)

		s.Equal(s.author(), persistedAuthor)
	})

	s.Run("It updates author when it already exists", func() {
		s.repo.SaveAuthor(s.ctx, s.author())

		author, err := s.usecase.Run(s.ctx, s.updatedInput())

		s.Equal(s.updatedAuthor(), author)
		s.Nil(err)

		persistedAuthor, _ := s.repo.GetAuthorByID(s.ctx, s.author().ID)

		s.Equal(s.updatedAuthor(), persistedAuthor)
	})
}

func (s *SaveAuthorUseCaseSuite) input() discussion.SaveAuthorInput {
	return discussion.SaveAuthorInput{
		ID:        "ID",
		Name:      "Name",
		AvatarURL: "https://example.com/avatar",
	}
}

func (s *SaveAuthorUseCaseSuite) updatedInput() discussion.SaveAuthorInput {
	return discussion.SaveAuthorInput{
		ID:        s.input().ID,
		Name:      "Updated Name",
		AvatarURL: "https://example.com/updated-avatar",
	}
}

func (s *SaveAuthorUseCaseSuite) author() *discussion.Author {
	input := s.input()

	return &discussion.Author{
		ID:        input.ID,
		Name:      input.Name,
		AvatarURL: input.AvatarURL,
	}
}

func (s *SaveAuthorUseCaseSuite) updatedAuthor() *discussion.Author {
	input := s.updatedInput()

	return &discussion.Author{
		ID:        input.ID,
		Name:      input.Name,
		AvatarURL: input.AvatarURL,
	}
}

func TestSaveAuthorUseCaseSuite(t *testing.T) {
	suite.Run(t, new(SaveAuthorUseCaseSuite))
}