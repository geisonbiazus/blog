package cache

import "github.com/geisonbiazus/blog/internal/adapters/cache/memory"

func NewMemoryCache() *memory.Cache {
	return memory.NewCache()
}
