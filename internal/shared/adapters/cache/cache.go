package cache

import (
	"github.com/geisonbiazus/blog/internal/shared/adapters/cache/memory"
	"github.com/geisonbiazus/blog/internal/shared/adapters/cache/null"
)

func NewMemoryCache() *memory.Cache {
	return memory.NewCache()
}

func NewNullCache() *null.Cache {
	return null.NewCache()
}
