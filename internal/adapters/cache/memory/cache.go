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

func (c *Cache[T]) Do(key string, resolve func() T, expiresIn time.Duration) T {
	item, ok := c.storage[key]
	if !ok || item.isExpired() {
		item = newCachedItem(resolve(), expiresIn)
		c.storage[key] = item
	}
	return item.value
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
