package null

import (
	"time"

	"github.com/geisonbiazus/blog/pkg/caching/ports"
)

type Cache struct{}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Do(
	key string,
	resolve ports.ResolveFn,
	expiresIn time.Duration,
) (interface{}, error) {
	return resolve()
}
