package discussion_test

import (
	"context"
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/commentrepo/memory"
	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestComment(t *testing.T) {
	t.Run("Replies", func(t *testing.T) {
		t.Run("It returns comment replies when they exist", func(t *testing.T) {
			ctx := context.Background()
			repo := memory.NewCommentRepo()

			comment := discussion.NewComment(discussion.CommentParams{
				ID: "comment",
			}, repo)

			reply1 := discussion.NewComment(discussion.CommentParams{
				ID:        "reply1",
				SubjectID: "comment",
			}, repo)

			reply2 := discussion.NewComment(discussion.CommentParams{
				ID:        "reply2",
				SubjectID: "comment",
			}, repo)

			repo.Save(ctx, comment)
			repo.Save(ctx, reply1)
			repo.Save(ctx, reply2)

			replies, err := comment.Replies(ctx)

			assert.DeepEqual(t, []*discussion.Comment{reply1, reply2}, replies)
			assert.Nil(t, err)
		})
	})
}
