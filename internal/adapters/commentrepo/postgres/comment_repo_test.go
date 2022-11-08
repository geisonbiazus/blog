package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/adapters/commentrepo/postgres"
	"github.com/geisonbiazus/blog/internal/adapters/idgenerator/uuid"
	"github.com/geisonbiazus/blog/internal/core/discussion"
	. "github.com/geisonbiazus/blog/internal/core/discussion/test"
	"github.com/geisonbiazus/blog/pkg/assert"
	"github.com/geisonbiazus/blog/pkg/dbrepo"
)

func TestCommentRepo(t *testing.T) {
	t.Run("GetCommentsAndRepliesRecursively", func(t *testing.T) {
		t.Run("It fetches the comments by subjectID", func(t *testing.T) {
			dbrepo.Test(func(ctx context.Context, db *sql.DB) {
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

				repo := postgres.NewCommentRepo(db)

				assert.Nil(t, repo.SaveAuthor(ctx, author))
				assert.Nil(t, repo.SaveComment(ctx, comment1))
				assert.Nil(t, repo.SaveComment(ctx, comment2))

				comments, err := repo.GetCommentsAndRepliesRecursively(ctx, subjectID)

				assert.Nil(t, err)
				assert.DeepEqual(t, []*discussion.Comment{
					comment1,
					comment2,
				}, comments)

				fmt.Println(comment1)
				fmt.Println(comment2)

				for _, c := range comments {
					fmt.Println(c)
				}
			})
		})
	})
}
