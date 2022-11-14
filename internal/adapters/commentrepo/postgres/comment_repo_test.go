package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/adapters/commentrepo/postgres"
	"github.com/geisonbiazus/blog/internal/adapters/idgenerator/uuid"
	"github.com/geisonbiazus/blog/internal/core/discussion"
	. "github.com/geisonbiazus/blog/internal/core/discussion/test"
	"github.com/geisonbiazus/blog/pkg/dbrepo"
	"github.com/stretchr/testify/assert"
)

type commentRepoFixture struct {
	repo                *postgres.CommentRepo
	uuidGen             *uuid.Generator
	author              *discussion.Author
	subjectID           string
	comment1            *discussion.Comment
	comment2            *discussion.Comment
	reply1              *discussion.Comment
	reply2              *discussion.Comment
	comment1WithReplies *discussion.Comment
}

func TestCommentRepo(t *testing.T) {
	setup := func(db *sql.DB) *commentRepoFixture {
		uuidGen := uuid.NewGenerator()

		author := NewAuthor(discussion.Author{ID: uuidGen.Generate()})
		subjectID := "SUBJECT_ID"

		comment1 := NewComment(discussion.Comment{
			ID:        uuidGen.Generate(),
			SubjectID: subjectID,
			AuthorID:  author.ID,
			Author:    author,
			CreatedAt: time.Date(2022, time.October, 4, 9, 0, 0, 0, time.UTC),
		})

		comment2 := NewComment(discussion.Comment{
			ID:        uuidGen.Generate(),
			SubjectID: subjectID,
			AuthorID:  author.ID,
			Author:    author,
			CreatedAt: time.Date(2022, time.October, 4, 10, 0, 0, 0, time.UTC),
		})

		reply1 := NewComment(discussion.Comment{
			ID:        uuidGen.Generate(),
			SubjectID: comment1.ID,
			AuthorID:  author.ID,
			Author:    author,
			CreatedAt: time.Date(2022, time.October, 4, 11, 0, 0, 0, time.UTC),
		})

		reply2 := NewComment(discussion.Comment{
			ID:        uuidGen.Generate(),
			SubjectID: reply1.ID,
			AuthorID:  author.ID,
			Author:    author,
			CreatedAt: time.Date(2022, time.October, 4, 12, 0, 0, 0, time.UTC),
		})

		reply1WithReplies := NewComment(*reply1)
		reply1WithReplies.Replies = []*discussion.Comment{reply2}

		comment1WithReplies := NewComment(*comment1)
		comment1WithReplies.Replies = []*discussion.Comment{reply1WithReplies}

		repo := postgres.NewCommentRepo(db)

		return &commentRepoFixture{
			repo:                repo,
			uuidGen:             uuidGen,
			author:              author,
			subjectID:           subjectID,
			comment1:            comment1,
			comment2:            comment2,
			reply1:              reply1,
			reply2:              reply2,
			comment1WithReplies: comment1WithReplies,
		}
	}

	t.Run("GetCommentsAndRepliesRecursively", func(t *testing.T) {
		t.Run("It fetches the comments by subjectID", func(t *testing.T) {
			dbrepo.Test(func(ctx context.Context, db *sql.DB) {
				f := setup(db)

				assert.Nil(t, f.repo.SaveAuthor(ctx, f.author))
				assert.Nil(t, f.repo.SaveComment(ctx, f.comment1))
				assert.Nil(t, f.repo.SaveComment(ctx, f.comment2))

				comments, err := f.repo.GetCommentsAndRepliesRecursively(ctx, f.subjectID)

				assert.Nil(t, err)
				assert.Equal(t, []*discussion.Comment{
					f.comment1,
					f.comment2,
				}, comments)
			})
		})

		t.Run("It fetches replies recursively", func(t *testing.T) {
			dbrepo.Test(func(ctx context.Context, db *sql.DB) {
				f := setup(db)

				assert.Nil(t, f.repo.SaveAuthor(ctx, f.author))
				assert.Nil(t, f.repo.SaveComment(ctx, f.comment1))
				assert.Nil(t, f.repo.SaveComment(ctx, f.reply1))
				assert.Nil(t, f.repo.SaveComment(ctx, f.reply2))

				comments, err := f.repo.GetCommentsAndRepliesRecursively(ctx, f.subjectID)

				assert.Nil(t, err)
				assert.Equal(t, []*discussion.Comment{
					f.comment1WithReplies,
				}, comments)
			})
		})
	})
}
