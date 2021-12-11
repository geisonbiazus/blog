package staterepo

import "github.com/geisonbiazus/blog/internal/adapters/staterepo/memory"

func NewMemoryStateRepo() *memory.StateRepo {
	return memory.NewStateRepo()
}
