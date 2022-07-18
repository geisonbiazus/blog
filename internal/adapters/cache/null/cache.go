package null

import (
	"time"

	"github.com/geisonbiazus/blog/internal/core/shared"
)

type Cache struct{}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Do(
	key string, resolve shared.ResolveFn, expiresIn time.Duration,
) (interface{}, error) {
	return resolve()
}
