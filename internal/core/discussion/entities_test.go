package discussion_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/core/discussion"
	. "github.com/geisonbiazus/blog/internal/core/discussion/test"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestComment(t *testing.T) {
	t.Run("Clone", func(t *testing.T) {
		t.Run("It makes a copy of the comment", func(t *testing.T) {
			comment := NewComment(discussion.Comment{})
			clone := comment.Clone()

			assert.False(t, comment == clone)
		})
	})
}
