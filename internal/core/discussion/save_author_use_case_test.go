package discussion_test

import (
	"context"
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/commentrepo/memory"
	"github.com/geisonbiazus/blog/internal/adapters/idgenerator/fake"
	"github.com/geisonbiazus/blog/internal/adapters/transactionmanager"
	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/stretchr/testify/suite"
)

type SaveAuthorUseCaseSuite struct {
	suite.Suite
	usecase *discussion.SaveAuthorUseCase
	repo    *memory.CommentRepo
	ctx     context.Context
	idGen   *fake.IDGenerator
}

func (s *SaveAuthorUseCaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = memory.NewCommentRepo()
	txManager := transactionmanager.NewFakeTransactionManager()
	s.idGen = fake.NewIDGenerator()
	s.idGen.ReturnID = s.AuthorID()
	s.usecase = discussion.NewSaveAuthorUseCase(s.repo, txManager, s.idGen)
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

func (s *SaveAuthorUseCaseSuite) input() discussion.SaveAuthorInput {
	return discussion.SaveAuthorInput{
		UserID:    "USER_ID",
		Name:      "Name",
		AvatarURL: "https://example.com/avatar",
	}
}

func (s *SaveAuthorUseCaseSuite) updatedInput() discussion.SaveAuthorInput {
	return discussion.SaveAuthorInput{
		UserID:    s.input().UserID,
		Name:      "Updated Name",
		AvatarURL: "https://example.com/updated-avatar",
	}
}

func (s *SaveAuthorUseCaseSuite) author() *discussion.Author {
	input := s.input()

	return &discussion.Author{
		ID:        s.AuthorID(),
		UserID:    input.UserID,
		Name:      input.Name,
		AvatarURL: input.AvatarURL,
	}
}

func (s *SaveAuthorUseCaseSuite) updatedAuthor() *discussion.Author {
	input := s.updatedInput()

	return &discussion.Author{
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
