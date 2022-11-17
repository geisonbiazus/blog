package discussion_test

import (
	"context"
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/adapters/commentrepo/memory"
	"github.com/geisonbiazus/blog/internal/core/discussion"
	. "github.com/geisonbiazus/blog/internal/core/discussion/test"
	"github.com/stretchr/testify/assert"
)

type listCommentsUseCaseFixture struct {
	usecase *discussion.ListCommentsUseCase
	repo    *memory.CommentRepo
	ctx     context.Context
	author  *discussion.Author
}

func TestListCommentsUseCase(t *testing.T) {
	setup := func() *listCommentsUseCaseFixture {
		ctx := context.Background()
		repo := memory.NewCommentRepo()
		usecase := discussion.NewListCommentsUseCase(repo)
		author := &discussion.Author{
			ID:        "AUTHOR_ID",
			Name:      "Author",
			AvatarURL: "https://example.com/avatar",
		}

		repo.SaveAuthor(ctx, author)

		return &listCommentsUseCaseFixture{
			usecase: usecase,
			repo:    repo,
			ctx:     ctx,
			author:  author,
		}
	}

	t.Run("It returns an empty list when there is no commments", func(t *testing.T) {
		f := setup()
		subjectID := "SUBJECT_ID"

		result, err := f.usecase.Run(f.ctx, subjectID)

		assert.Equal(t, []*discussion.Comment{}, result)
		assert.Nil(t, err)
	})

	t.Run("It fetches and returns the comments of the given subject in chronological order", func(t *testing.T) {
		f := setup()

		comment1 := NewComment(discussion.Comment{
			ID:        "ID_1",
			AuthorID:  f.author.ID,
			Author:    f.author,
			CreatedAt: time.Date(2022, time.October, 4, 9, 0, 0, 0, time.UTC),
		})

		comment2 := NewComment(discussion.Comment{
			ID:        "ID_2",
			AuthorID:  f.author.ID,
			Author:    f.author,
			CreatedAt: time.Date(2022, time.October, 4, 8, 0, 0, 0, time.UTC),
		})

		f.repo.SaveComment(f.ctx, comment1)
		f.repo.SaveComment(f.ctx, comment2)

		result, err := f.usecase.Run(f.ctx, comment1.SubjectID)

		assert.Equal(t, []*discussion.Comment{comment2, comment1}, result)
		assert.Nil(t, err)
	})

	t.Run("It fetches replies recursively", func(t *testing.T) {
		f := setup()

		comment := NewComment(discussion.Comment{
			ID:       "COMMENT",
			AuthorID: f.author.ID,
		})

		reply1 := NewComment(discussion.Comment{
			ID:        "REPLY_1",
			SubjectID: comment.ID,
			AuthorID:  f.author.ID,
		})

		reply2 := NewComment(discussion.Comment{
			ID:        "REPLY_2",
			SubjectID: reply1.ID,
			AuthorID:  f.author.ID,
		})

		f.repo.SaveComment(f.ctx, comment)
		f.repo.SaveComment(f.ctx, reply1)
		f.repo.SaveComment(f.ctx, reply2)

		result, err := f.usecase.Run(f.ctx, comment.SubjectID)

		// TODO: Return author

		commentWithReplies := []*discussion.Comment{
			NewComment(discussion.Comment{
				ID:        comment.ID,
				CreatedAt: comment.CreatedAt,
				AuthorID:  comment.AuthorID,
				Author:    f.author,
				Replies: []*discussion.Comment{
					NewComment(discussion.Comment{
						ID:        reply1.ID,
						SubjectID: reply1.SubjectID,
						CreatedAt: reply1.CreatedAt,
						AuthorID:  reply1.AuthorID,
						Author:    f.author,
						Replies: []*discussion.Comment{
							NewComment(discussion.Comment{
								ID:        reply2.ID,
								SubjectID: reply2.SubjectID,
								CreatedAt: reply2.CreatedAt,
								AuthorID:  reply2.AuthorID,
								Author:    f.author,
							}),
						},
					}),
				},
			}),
		}

		assert.Equal(t, commentWithReplies, result)
		assert.Equal(t, commentWithReplies[0].Replies[0].Replies[0], result[0].Replies[0].Replies[0])
		assert.Nil(t, err)
	})
}
