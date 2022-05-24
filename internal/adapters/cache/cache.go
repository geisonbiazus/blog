package cache

import "github.com/geisonbiazus/blog/internal/adapters/cache/memory"

func NewMemoryCache[T any]() *memory.Cache[T] {
	return memory.NewCache[T]()
}
