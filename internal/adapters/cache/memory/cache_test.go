package memory_test

import (
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/adapters/cache/memory"
	"github.com/geisonbiazus/blog/internal/core/shared"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestCache(t *testing.T) {
	t.Run("Do", func(t *testing.T) {
		t.Run("It executes the resolve fn and return its value first time it's called", func(t *testing.T) {
			cache := memory.NewCache[int]()
			calls := 0

			result := cache.Do("key", func() int {
				calls++
				return calls
			}, shared.NeverExpire)

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

			cache.Do("key", resolve, shared.NeverExpire)
			result := cache.Do("key", resolve, shared.NeverExpire)

			assert.Equal(t, 1, result)
			assert.Equal(t, 1, calls)
		})

		t.Run("It returns any value type", func(t *testing.T) {
			cache := memory.NewCache[string]()

			result := cache.Do("key", func() string {
				return "value"
			}, shared.NeverExpire)

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

			cache.Do("key1", resolve1, shared.NeverExpire)
			result1 := cache.Do("key1", resolve1, shared.NeverExpire)

			cache.Do("key2", resolve2, shared.NeverExpire)
			result2 := cache.Do("key2", resolve2, shared.NeverExpire)

			assert.Equal(t, "value1", result1)
			assert.Equal(t, "value2", result2)
			assert.Equal(t, 1, calls1)
			assert.Equal(t, 1, calls2)
		})

		t.Run("It expires the cache based on the given interval", func(t *testing.T) {
			cache := memory.NewCache[int]()
			calls := 0

			resolve := func() int {
				calls++
				return calls
			}

			cache.Do("key", resolve, -1*time.Minute)
			cache.Do("key", resolve, -1*time.Minute)

			assert.Equal(t, 2, calls)
		})
	})
}
