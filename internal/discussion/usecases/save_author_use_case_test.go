package usecases_test

import (
	"context"
	"testing"

	"github.com/geisonbiazus/blog/internal/discussion/adapters/commentrepo/memory"
	"github.com/geisonbiazus/blog/internal/discussion/entities"
	"github.com/geisonbiazus/blog/internal/discussion/usecases"
	"github.com/geisonbiazus/blog/pkg/gen"
	"github.com/geisonbiazus/blog/pkg/gen/fake"
	"github.com/geisonbiazus/blog/pkg/transaction"
	"github.com/stretchr/testify/suite"
)

type SaveAuthorUseCaseSuite struct {
	suite.Suite
	usecase *usecases.SaveAuthorUseCase
	repo    *memory.CommentRepo
	ctx     context.Context
	idGen   *fake.Generator
}

func (s *SaveAuthorUseCaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = memory.NewCommentRepo()
	txManager := transaction.NewFakeManager()
	s.idGen = gen.NewFakeGenerator()
	s.idGen.ReturnID = s.AuthorID()
	s.usecase = usecases.NewSaveAuthorUseCase(s.repo, txManager, s.idGen)
}

func (s *SaveAuthorUseCaseSuite) TestRun() {
	s.Run("It creates an author when it doesn't exist", func() {
		author, err := s.usecase.Run(s.ctx, s.input())

		s.Equal(s.author(), author)
		s.Nil(err)

		persistedAuthor, _ := s.repo.GetAuthorByID(s.ctx, author.ID)

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

func (s *SaveAuthorUseCaseSuite) input() usecases.SaveAuthorInput {
	return usecases.SaveAuthorInput{
		UserID:    "USER_ID",
		Name:      "Name",
		AvatarURL: "https://example.com/avatar",
	}
}

func (s *SaveAuthorUseCaseSuite) updatedInput() usecases.SaveAuthorInput {
	return usecases.SaveAuthorInput{
		UserID:    s.input().UserID,
		Name:      "Updated Name",
		AvatarURL: "https://example.com/updated-avatar",
	}
}

func (s *SaveAuthorUseCaseSuite) author() *entities.Author {
	input := s.input()

	return &entities.Author{
		ID:        s.AuthorID(),
		UserID:    input.UserID,
		Name:      input.Name,
		AvatarURL: input.AvatarURL,
	}
}

func (s *SaveAuthorUseCaseSuite) updatedAuthor() *entities.Author {
	input := s.updatedInput()

	return &entities.Author{
		ID:        s.AuthorID(),
		UserID:    input.UserID,
		Name:      input.Name,
		AvatarURL: input.AvatarURL,
	}
}

func (s *SaveAuthorUseCaseSuite) AuthorID() string {
	return "AUTHOR_ID"
}

func TestSaveAuthorUseCaseSuite(t *testing.T) {
	suite.Run(t, new(SaveAuthorUseCaseSuite))
}
