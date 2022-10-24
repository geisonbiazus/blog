package discussion_test

import (
	"context"
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/adapters/commentrepo/memory"
	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type listCommentsUseCaseFixture struct {
	usecase *discussion.ListCommentsUseCase
	repo    *memory.CommentRepo
	ctx     context.Context
}

func TestListCommentsUseCase(t *testing.T) {
	setup := func() *listCommentsUseCaseFixture {
		repo := memory.NewCommentRepo()
		usecase := discussion.NewListCommentsUseCase(repo)

		return &listCommentsUseCaseFixture{
			usecase: usecase,
			repo:    repo,
			ctx:     context.Background(),
		}
	}

	t.Run("It returns an empty list when there is no commments", func(t *testing.T) {
		f := setup()
		subjectID := "SUBJECT_ID"

		result, err := f.usecase.Run(f.ctx, subjectID)

		assert.DeepEqual(t, []*discussion.Comment{}, result)
		assert.Nil(t, err)
	})

	t.Run("It fetches and returns the comments of the given subject in chronological order", func(t *testing.T) {
		f := setup()

		comment1 := newComment(discussion.Comment{
			ID:        "ID_1",
			CreatedAt: time.Date(2022, time.October, 4, 9, 0, 0, 0, time.UTC),
		})

		comment2 := newComment(discussion.Comment{
			ID:        "ID_2",
			CreatedAt: time.Date(2022, time.October, 4, 8, 0, 0, 0, time.UTC),
		})

		f.repo.SaveComment(f.ctx, comment1)
		f.repo.SaveComment(f.ctx, comment2)

		result, err := f.usecase.Run(f.ctx, comment1.SubjectID)

		assert.DeepEqual(t, []*discussion.Comment{comment2, comment1}, result)
		assert.Nil(t, err)
	})

	t.Run("It fetches replies recursively", func(t *testing.T) {
		f := setup()

		comment := newComment(discussion.Comment{
			ID: "COMMENT",
		})

		reply1 := newComment(discussion.Comment{
			ID:        "REPLY_1",
			SubjectID: comment.ID,
		})

		reply2 := newComment(discussion.Comment{
			ID:        "REPLY_2",
			SubjectID: reply1.ID,
		})

		f.repo.SaveComment(f.ctx, comment)
		f.repo.SaveComment(f.ctx, reply1)
		f.repo.SaveComment(f.ctx, reply2)

		result, err := f.usecase.Run(f.ctx, comment.SubjectID)

		commentWithReplies := []*discussion.Comment{
			newComment(discussion.Comment{
				ID:        comment.ID,
				CreatedAt: comment.CreatedAt,
				Replies: []*discussion.Comment{
					newComment(discussion.Comment{
						ID:        reply1.ID,
						SubjectID: reply1.SubjectID,
						CreatedAt: reply1.CreatedAt,
						Replies: []*discussion.Comment{
							newComment(discussion.Comment{
								ID:        reply2.ID,
								SubjectID: reply2.SubjectID,
								CreatedAt: reply2.CreatedAt,
							}),
						},
					}),
				},
			}),
		}

		assert.DeepEqual(t, commentWithReplies, result)
		assert.DeepEqual(t, commentWithReplies[0].Replies[0].Replies[0], result[0].Replies[0].Replies[0])
		assert.Nil(t, err)
	})
}
