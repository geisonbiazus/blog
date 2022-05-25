package memory

import (
	"time"

	"github.com/geisonbiazus/blog/internal/core/shared"
)

type Cache struct {
	storage map[string]cachedItem
}

func NewCache() *Cache {
	return &Cache{
		storage: map[string]cachedItem{},
	}
}

func (c *Cache) Do(
	key string,
	resolve shared.ResolveFn,
	expiresIn time.Duration,
) (interface{}, error) {
	item, ok := c.storage[key]
	if !ok || item.isExpired() {
		return c.resolveAndStoreValue(key, resolve, expiresIn)
	}
	return item.value, nil
}

func (c *Cache) resolveAndStoreValue(
	key string,
	resolve shared.ResolveFn,
	expiresIn time.Duration,
) (interface{}, error) {
	value, err := resolve()

	if err == nil {
		c.storage[key] = newCachedItem(value, expiresIn)
	}

	return value, err
}

type cachedItem struct {
	value     interface{}
	createdAt time.Time
	expiresIn time.Duration
}

func newCachedItem(value interface{}, expiresIn time.Duration) cachedItem {
	return cachedItem{
		value:     value,
		createdAt: time.Now(),
		expiresIn: expiresIn,
	}
}

func (i cachedItem) isExpired() bool {
	if i.expiresIn == shared.NeverExpire {
		return false
	}

	return i.expiresAt().Before(time.Now())
}

func (i cachedItem) expiresAt() time.Time {
	return i.createdAt.Add(i.expiresIn)
}
