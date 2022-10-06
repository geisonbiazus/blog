package discussion_test

import (
	"context"
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/adapters/commentrepo/memory"
	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type listCommentFixture struct {
	comment *discussion.Comment
	reply1  *discussion.Comment
	reply2  *discussion.Comment
	repo    *memory.CommentRepo
	ctx     context.Context
}

func TestComment(t *testing.T) {
	setup := func() *listCommentFixture {
		ctx := context.Background()
		repo := memory.NewCommentRepo()
		comment := newComment(discussion.CommentParams{}, repo)
		reply1 := newComment(discussion.CommentParams{
			ID:        "reply1",
			SubjectID: comment.ID,
			CreatedAt: time.Now().Add(-1 * time.Hour),
		}, repo)

		reply2 := newComment(discussion.CommentParams{
			ID:        "reply2",
			SubjectID: comment.ID,
			CreatedAt: time.Now(),
		}, repo)

		repo.Save(ctx, comment)
		repo.Save(ctx, reply1)
		repo.Save(ctx, reply2)

		return &listCommentFixture{
			comment: comment,
			reply1:  reply1,
			reply2:  reply2,
			repo:    repo,
			ctx:     ctx,
		}
	}

	t.Run("NewComment", func(t *testing.T) {
		t.Run("it initializes a comment with all params", func(t *testing.T) {
			f := setup()

			params := discussion.CommentParams{
				ID:        "ID",
				SubjectID: "SUBJECT_ID",
				Markdown:  "Markdown",
				HTML:      "HTML",
				CreatedAt: time.Now(),
			}

			comment := discussion.NewComment(params, f.repo)

			assert.Equal(t, params.ID, comment.ID)
			assert.Equal(t, params.SubjectID, comment.SubjectID)
			assert.Equal(t, params.Markdown, comment.Markdown)
			assert.Equal(t, params.HTML, comment.HTML)
			assert.Equal(t, params.CreatedAt, comment.CreatedAt)
		})
	})

	t.Run("Replies", func(t *testing.T) {
		t.Run("It returns comment replies when they exist", func(t *testing.T) {
			f := setup()

			replies, err := f.comment.Replies(f.ctx)

			assert.DeepEqual(t, []*discussion.Comment{f.reply1, f.reply2}, replies)
			assert.Nil(t, err)
		})

		t.Run("It caches loaded replies", func(t *testing.T) {
			f := setup()

			f.comment.Replies(f.ctx)

			reply3 := newComment(discussion.CommentParams{ID: "repply3", SubjectID: f.comment.ID}, f.repo)

			f.repo.Save(f.ctx, reply3)

			replies, err := f.comment.Replies(f.ctx)

			assert.DeepEqual(t, []*discussion.Comment{f.reply1, f.reply2}, replies)
			assert.Nil(t, err)
		})
	})

	t.Run("SetReplies", func(t *testing.T) {
		t.Run("It sets the replies so they are not loaded when requested", func(t *testing.T) {
			f := setup()

			f.comment.SetReplies([]*discussion.Comment{f.reply2})

			replies, err := f.comment.Replies(f.ctx)

			assert.DeepEqual(t, []*discussion.Comment{f.reply2}, replies)
			assert.Nil(t, err)
		})
	})
}
