package entities_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/discussion/entities"
	. "github.com/geisonbiazus/blog/internal/discussion/test"
	"github.com/stretchr/testify/assert"
)

func TestComment(t *testing.T) {
	t.Run("Clone", func(t *testing.T) {
		t.Run("It makes a copy of the comment", func(t *testing.T) {
			comment := NewComment(entities.Comment{})
			clone := comment.Clone()

			assert.False(t, comment == clone)
		})
	})
}
