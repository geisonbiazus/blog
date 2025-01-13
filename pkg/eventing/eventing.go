package eventing

import (
	"github.com/geisonbiazus/blog/pkg/eventing/entities"
	"github.com/geisonbiazus/blog/pkg/eventing/fake"
	"github.com/geisonbiazus/blog/pkg/eventing/memory"
	"github.com/geisonbiazus/blog/pkg/eventing/ports"
)

type Event = entities.Event
type Publisher = ports.Publisher

func NewFakePublisher() *fake.Publisher {
	return fake.NewPublisher()
}

func NewMemoryPubSub() *memory.PubSub {
	return memory.NewPubSub()
}
