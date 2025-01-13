package caching

import (
	"github.com/geisonbiazus/blog/pkg/caching/memory"
	"github.com/geisonbiazus/blog/pkg/caching/null"
	"github.com/geisonbiazus/blog/pkg/caching/ports"
)

func NewMemoryCache() *memory.Cache {
	return memory.NewCache()
}

func NewNullCache() *null.Cache {
	return null.NewCache()
}

type ResolveFn = ports.ResolveFn
type Cache = ports.Cache

var NeverExpire = ports.NeverExpire
