package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/discussion/adapters/commentrepo/postgres"
	"github.com/geisonbiazus/blog/internal/discussion/entities"
	. "github.com/geisonbiazus/blog/internal/discussion/test"
	"github.com/geisonbiazus/blog/pkg/dbrepo"
	"github.com/geisonbiazus/blog/pkg/gen/uuid"
	"github.com/stretchr/testify/suite"
)

type CommentRepoSuite struct {
	suite.Suite
	repo                *postgres.CommentRepo
	uuidGen             *uuid.Generator
	author              *entities.Author
	subjectID           string
	comment1            *entities.Comment
	comment2            *entities.Comment
	reply1              *entities.Comment
	reply2              *entities.Comment
	comment1WithReplies *entities.Comment
}

func (s *CommentRepoSuite) SetupSubTest() {
	s.uuidGen = uuid.NewGenerator()

	s.author = NewAuthor(entities.Author{ID: s.uuidGen.Generate(), UserID: s.uuidGen.Generate()})
	s.subjectID = "SUBJECT_ID"

	s.comment1 = NewComment(entities.Comment{
		ID:        s.uuidGen.Generate(),
		SubjectID: s.subjectID,
		AuthorID:  s.author.ID,
		Author:    s.author,
		CreatedAt: time.Date(2022, time.October, 4, 9, 0, 0, 0, time.UTC),
	})

	s.comment2 = NewComment(entities.Comment{
		ID:        s.uuidGen.Generate(),
		SubjectID: s.subjectID,
		AuthorID:  s.author.ID,
		Author:    s.author,
		CreatedAt: time.Date(2022, time.October, 4, 10, 0, 0, 0, time.UTC),
	})

	s.reply1 = NewComment(entities.Comment{
		ID:        s.uuidGen.Generate(),
		SubjectID: s.comment1.ID,
		AuthorID:  s.author.ID,
		Author:    s.author,
		CreatedAt: time.Date(2022, time.October, 4, 11, 0, 0, 0, time.UTC),
	})

	s.reply2 = NewComment(entities.Comment{
		ID:        s.uuidGen.Generate(),
		SubjectID: s.reply1.ID,
		AuthorID:  s.author.ID,
		Author:    s.author,
		CreatedAt: time.Date(2022, time.October, 4, 12, 0, 0, 0, time.UTC),
	})

	reply1WithReplies := NewComment(*s.reply1)
	reply1WithReplies.Replies = []*entities.Comment{s.reply2}

	s.comment1WithReplies = NewComment(*s.comment1)
	s.comment1WithReplies.Replies = []*entities.Comment{reply1WithReplies}
}

func (s *CommentRepoSuite) TestGetAuthorByID() {
	s.Run("When author doesn't exist", func() {
		s.Run("It returns nil", func() {
			dbrepo.Test(func(ctx context.Context, db *sql.DB) {
				repo := postgres.NewCommentRepo(db)
				author, err := repo.GetAuthorByID(ctx, s.uuidGen.Generate())

				s.Nil(err)
				s.Nil(author)
			})
		})

		s.Run("When author exists", func() {
			s.Run("It returns the author", func() {
				dbrepo.Test(func(ctx context.Context, db *sql.DB) {
					repo := postgres.NewCommentRepo(db)

					err := repo.SaveAuthor(ctx, s.author)
					s.Nil(err)

					author, err := repo.GetAuthorByID(ctx, s.author.ID)

					s.Nil(err)
					s.Equal(s.author, author)
					s.True(author.Persisted)
				})
			})
		})
	})
}

func (s *CommentRepoSuite) TestGetAuthorByUserId() {
	s.Run("When author doesn't exist", func() {
		s.Run("It returns nil", func() {
			dbrepo.Test(func(ctx context.Context, db *sql.DB) {
				repo := postgres.NewCommentRepo(db)
				author, err := repo.GetAuthorByUserID(ctx, s.uuidGen.Generate())

				s.Nil(err)
				s.Nil(author)
			})
		})
	})

	s.Run("When author exists", func() {
		s.Run("It returns the author", func() {
			dbrepo.Test(func(ctx context.Context, db *sql.DB) {
				repo := postgres.NewCommentRepo(db)

				err := repo.SaveAuthor(ctx, s.author)
				s.Nil(err)

				author, err := repo.GetAuthorByUserID(ctx, s.author.UserID)

				s.Nil(err)
				s.Equal(s.author, author)
				s.True(author.Persisted)
			})
		})
	})
}

func (s *CommentRepoSuite) TestSaveAuthor() {
	s.Run("With a new author", func() {
		s.Run("It inserts the author", func() {
			dbrepo.Test(func(ctx context.Context, db *sql.DB) {
				repo := postgres.NewCommentRepo(db)

				err := repo.SaveAuthor(ctx, s.author)
				s.Nil(err)
				s.True(s.author.Persisted)

				author, err := repo.GetAuthorByID(ctx, s.author.ID)

				s.Nil(err)
				s.Equal(s.author, author)
			})
		})
	})

	s.Run("With an existing author", func() {
		s.Run("It updates the author", func() {
			dbrepo.Test(func(ctx context.Context, db *sql.DB) {
				repo := postgres.NewCommentRepo(db)

				repo.SaveAuthor(ctx, s.author)

				author, _ := repo.GetAuthorByID(ctx, s.author.ID)

				author.Name = "Updated Name"
				author.AvatarURL = "https://example.com/updated-avatar"
				author.UserID = s.uuidGen.Generate()

				err := repo.SaveAuthor(ctx, author)
				s.Nil(err)

				updatedAuthor, err := repo.GetAuthorByID(ctx, s.author.ID)

				s.Nil(err)
				s.Equal(author, updatedAuthor)
			})
		})
	})
}

func (s *CommentRepoSuite) TestGetCommentsAndRepliesRecursively() {
	s.Run("It fetches the comments by subjectID", func() {
		dbrepo.Test(func(ctx context.Context, db *sql.DB) {
			s.repo = postgres.NewCommentRepo(db)

			s.Nil(s.repo.SaveAuthor(ctx, s.author))
			s.Nil(s.repo.SaveComment(ctx, s.comment1))
			s.Nil(s.repo.SaveComment(ctx, s.comment2))

			comments, err := s.repo.GetCommentsAndRepliesRecursively(ctx, s.subjectID)

			s.Nil(err)
			s.Equal([]*entities.Comment{
				s.comment1,
				s.comment2,
			}, comments)
		})
	})

	s.Run("It fetches replies recursively", func() {
		dbrepo.Test(func(ctx context.Context, db *sql.DB) {
			s.repo = postgres.NewCommentRepo(db)

			s.Nil(s.repo.SaveAuthor(ctx, s.author))
			s.Nil(s.repo.SaveComment(ctx, s.comment1))
			s.Nil(s.repo.SaveComment(ctx, s.reply1))
			s.Nil(s.repo.SaveComment(ctx, s.reply2))

			comments, err := s.repo.GetCommentsAndRepliesRecursively(ctx, s.subjectID)

			s.Nil(err)
			s.Equal([]*entities.Comment{
				s.comment1WithReplies,
			}, comments)
		})
	})
}

func TestCommentRepoSuite(t *testing.T) {
	suite.Run(t, new(CommentRepoSuite))
}
