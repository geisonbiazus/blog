package memory_test

import (
	"errors"
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/adapters/cache/memory"
	"github.com/geisonbiazus/blog/internal/core/shared"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestCache(t *testing.T) {
	t.Run("Do", func(t *testing.T) {
		t.Run("It executes the resolve fn and return its value", func(t *testing.T) {
			cache := memory.NewCache()

			result, err := cache.Do("key", func() (interface{}, error) {
				return 1, nil
			}, shared.NeverExpire)

			assert.Equal(t, 1, result)
			assert.Nil(t, err)
		})

		t.Run("It returns error when resolve fn returns it", func(t *testing.T) {
			cache := memory.NewCache()
			err := errors.New("error")

			_, retunedErr := cache.Do("key", func() (interface{}, error) {
				return 1, err
			}, shared.NeverExpire)

			assert.Equal(t, err, retunedErr)
		})

		t.Run("It caches the value and returns it on consecutive calls", func(t *testing.T) {
			cache := memory.NewCache()
			calls := 0

			resolve := func() (interface{}, error) {
				calls++
				return calls, nil
			}

			cache.Do("key", resolve, shared.NeverExpire)
			result, err := cache.Do("key", resolve, shared.NeverExpire)

			assert.Equal(t, 1, result)
			assert.Equal(t, 1, calls)
			assert.Nil(t, err)
		})

		t.Run("It does not cache the value when an error is returned", func(t *testing.T) {
			cache := memory.NewCache()
			err := errors.New("error")
			calls := 0

			resolve := func() (interface{}, error) {
				calls++
				return calls, err
			}

			cache.Do("key", resolve, shared.NeverExpire)
			result, returnedErr := cache.Do("key", resolve, shared.NeverExpire)

			assert.Equal(t, 2, result)
			assert.Equal(t, 2, calls)
			assert.Equal(t, err, returnedErr)
		})

		t.Run("It returns any value type", func(t *testing.T) {
			cache := memory.NewCache()

			result, err := cache.Do("key", func() (interface{}, error) {
				return "value", nil
			}, shared.NeverExpire)

			assert.Equal(t, "value", result)
			assert.Nil(t, err)
		})

		t.Run("It caches different values independently based on the key", func(t *testing.T) {
			cache := memory.NewCache()
			calls1 := 0
			calls2 := 0
			resolve1 := func() (interface{}, error) {
				calls1++
				return "value1", nil
			}
			resolve2 := func() (interface{}, error) {
				calls2++
				return "value2", nil
			}

			cache.Do("key1", resolve1, shared.NeverExpire)
			result1, _ := cache.Do("key1", resolve1, shared.NeverExpire)

			cache.Do("key2", resolve2, shared.NeverExpire)
			result2, _ := cache.Do("key2", resolve2, shared.NeverExpire)

			assert.Equal(t, "value1", result1)
			assert.Equal(t, "value2", result2)
			assert.Equal(t, 1, calls1)
			assert.Equal(t, 1, calls2)
		})

		t.Run("It expires the cache based on the given interval", func(t *testing.T) {
			cache := memory.NewCache()
			calls := 0

			resolve := func() (interface{}, error) {
				calls++
				return calls, nil
			}

			cache.Do("key", resolve, -1*time.Minute)
			cache.Do("key", resolve, -1*time.Minute)

			assert.Equal(t, 2, calls)
		})
	})
}
