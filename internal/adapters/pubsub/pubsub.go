package pubsub

import "github.com/geisonbiazus/blog/internal/adapters/pubsub/memory"

func NewMemoryPubSub() *memory.PubSub {
	return memory.NewPubSub()
}
