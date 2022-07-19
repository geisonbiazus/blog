package null_test

import (
	"errors"
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/cache/null"
	"github.com/geisonbiazus/blog/internal/core/shared"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestCache(t *testing.T) {
	t.Run("It intializes", func(t *testing.T) {
		var _ *null.Cache = null.NewCache()
	})

	t.Run("Do", func(t *testing.T) {
		t.Run("It returns the resolver value and error", func(t *testing.T) {
			cache := null.NewCache()
			value := "value"
			err := errors.New("error")

			returnedValue, returnErr := cache.Do("key", func() (interface{}, error) {
				return value, err
			}, shared.NeverExpire)

			assert.Equal(t, value, returnedValue)
			assert.Equal(t, err, returnErr)
		})

		t.Run("It doesn't cache any value", func(t *testing.T) {
			cache := null.NewCache()

			cache.Do("key", func() (interface{}, error) {
				return "first value", nil
			}, shared.NeverExpire)

			value, _ := cache.Do("key", func() (interface{}, error) {
				return "second value", nil
			}, shared.NeverExpire)

			assert.Equal(t, "second value", value)
		})
	})
}
