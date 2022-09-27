package discussion_test

import (
	"context"
	"testing"

	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type listCommentsUseCaseFixture struct {
	usecase *discussion.ListCommentsUseCase
	repo    *CommentRepo
	ctx     context.Context
}

func TestListCommentsUseCase(t *testing.T) {
	setup := func() *listCommentsUseCaseFixture {
		usecase := discussion.NewListCommentsUseCase()
		repo := NewCommentRepo()

		return &listCommentsUseCaseFixture{
			usecase: usecase,
			repo:    repo,
			ctx:     context.Background(),
		}
	}

	t.Run("It returns an empty list when there is no commments", func(t *testing.T) {
		f := setup()
		subjectID := "SUBJECT_ID"

		assert.DeepEqual(t, []discussion.RenderedComment{}, f.usecase.Run(subjectID))
	})

	t.Run("It fetches and returns the comments from the repository", func(t *testing.T) {
		f := setup()

		comment := discussion.Comment{
			ID:        "ID",
			SubjectID: "SUBJECT_ID",
			Body:      "Body",
		}

		f.repo.Save(f.ctx, comment)
	})
}

type CommentRepo struct{}

func NewCommentRepo() *CommentRepo {
	return &CommentRepo{}
}

func (r *CommentRepo) Save(ctx context.Context, comment discussion.Comment) error {
	return nil
}
