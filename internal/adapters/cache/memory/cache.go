package memory

type Cache[T any] struct {
	storage map[string]T
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		storage: map[string]T{},
	}
}

func (c *Cache[T]) With(key string, resolve func() T) T {
	value, ok := c.storage[key]
	if !ok {
		value = resolve()
		c.storage[key] = value
	}
	return value
}
