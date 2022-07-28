package playground

import (
	"testing"

	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestHashMap(t *testing.T) {
	t.Run("It initializes", func(t *testing.T) {
		var _ *HashMap[int] = NewHashMap[int]()
	})

	t.Run("It sets and gets a value", func(t *testing.T) {
		hashMap := NewHashMap[int]()
		hashMap.Set("key", 1)
		assert.Equal(t, 1, hashMap.Get("key"))
	})
}

type HashMap[T any] struct {
	buckets []T
}

func NewHashMap[T any]() *HashMap[T] {
	return &HashMap[T]{
		buckets: make([]T, 1000),
	}
}

func (h *HashMap[T]) Set(key string, value T) {
	h.buckets[0] = value
}

func (h *HashMap[T]) Get(key string) T {
	return h.buckets[0]
}
