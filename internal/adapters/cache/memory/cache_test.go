package memory_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/cache/memory"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestCache(t *testing.T) {
	t.Run("With", func(t *testing.T) {
		t.Run("It executes the resolve fn and return its value first time it's called", func(t *testing.T) {
			cache := memory.NewCache[int]()
			calls := 0

			result := cache.With("key", func() int {
				calls++
				return calls
			})

			assert.Equal(t, 1, result)
			assert.Equal(t, 1, calls)
		})

		t.Run("It caches the value and returns it on consecutive calls", func(t *testing.T) {
			cache := memory.NewCache[int]()
			calls := 0

			resolve := func() int {
				calls++
				return calls
			}

			cache.With("key", resolve)
			result := cache.With("key", resolve)

			assert.Equal(t, 1, result)
			assert.Equal(t, 1, calls)
		})

		t.Run("It returns any value type", func(t *testing.T) {
			cache := memory.NewCache[string]()

			result := cache.With("key", func() string {
				return "value"
			})

			assert.Equal(t, "value", result)
		})

		t.Run("It caches different values independently based on the key", func(t *testing.T) {
			cache := memory.NewCache[string]()
			calls1 := 0
			calls2 := 0
			resolve1 := func() string {
				calls1++
				return "value1"
			}
			resolve2 := func() string {
				calls2++
				return "value2"
			}

			cache.With("key1", resolve1)
			result1 := cache.With("key1", resolve1)

			cache.With("key2", resolve2)
			result2 := cache.With("key2", resolve2)

			assert.Equal(t, "value1", result1)
			assert.Equal(t, "value2", result2)
			assert.Equal(t, 1, calls1)
			assert.Equal(t, 1, calls2)
		})
	})
}
