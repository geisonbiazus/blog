package memory

import (
	"time"

	"github.com/geisonbiazus/blog/internal/core/shared"
)

type Cache[T any] struct {
	storage map[string]cachedItem[T]
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		storage: map[string]cachedItem[T]{},
	}
}

func (c *Cache[T]) Do(key string, resolve func() (T, error), expiresIn time.Duration) (T, error) {
	item, ok := c.storage[key]
	if !ok || item.isExpired() {
		return c.resolveAndStoreValue(key, resolve, expiresIn)
	}
	return item.value, nil
}

func (c *Cache[T]) resolveAndStoreValue(key string, resolve func() (T, error), expiresIn time.Duration) (T, error) {
	value, err := resolve()

	if err == nil {
		c.storage[key] = newCachedItem(value, expiresIn)
	}

	return value, err
}

type cachedItem[T any] struct {
	value     T
	createdAt time.Time
	expiresIn time.Duration
}

func newCachedItem[T any](value T, expiresIn time.Duration) cachedItem[T] {
	return cachedItem[T]{
		value:     value,
		createdAt: time.Now(),
		expiresIn: expiresIn,
	}
}

func (i cachedItem[T]) isExpired() bool {
	if i.expiresIn == shared.NeverExpire {
		return false
	}

	return i.expiresAt().Before(time.Now())
}

func (i cachedItem[T]) expiresAt() time.Time {
	return i.createdAt.Add(i.expiresIn)
}
