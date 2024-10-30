package staterepo

import "github.com/geisonbiazus/blog/internal/auth/adapters/staterepo/memory"

func NewMemoryStateRepo() *memory.StateRepo {
	return memory.NewStateRepo()
}
