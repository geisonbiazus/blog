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

		comment1 := newComment(discussion.CommentParams{
			ID:        "ID_1",
			CreatedAt: time.Date(2022, time.October, 4, 9, 0, 0, 0, time.UTC),
		}, f.repo)

		comment2 := newComment(discussion.CommentParams{
			ID:        "ID_2",
			CreatedAt: time.Date(2022, time.October, 4, 8, 0, 0, 0, time.UTC),
		}, f.repo)

		f.repo.SaveComment(f.ctx, comment1)
		f.repo.SaveComment(f.ctx, comment2)

		result, err := f.usecase.Run(f.ctx, comment1.SubjectID)

		assert.DeepEqual(t, []*discussion.Comment{comment2, comment1}, result)
		assert.Nil(t, err)
	})
}
