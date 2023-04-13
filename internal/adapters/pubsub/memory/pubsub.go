package memory

import (
	"github.com/geisonbiazus/blog/internal/core/shared"
)

type PubSub struct {
	SubscriberBuffer int
	subscribers      map[string][]chan shared.Event
}

func NewPubSub() *PubSub {
	return &PubSub{
		SubscriberBuffer: 10,
		subscribers:      map[string][]chan shared.Event{},
	}
}

func (p *PubSub) Publish(event shared.Event) error {
	for _, subscriber := range p.subscribers[event.Type] {
		subscriber <- event
	}

	return nil
}

func (p *PubSub) Subscribe(eventType string) chan shared.Event {
	channel := make(chan shared.Event, p.SubscriberBuffer)
	p.ensureSubscribersFor(eventType)
	p.subscribers[eventType] = append(p.subscribers[eventType], channel)

	return channel
}

func (p *PubSub) ensureSubscribersFor(eventType string) {
	if _, ok := p.subscribers[eventType]; !ok {
		p.subscribers[eventType] = []chan shared.Event{}
	}
}

func (p *PubSub) NotifyError(event shared.Event, err error) {}

func (p *PubSub) NotifySuccess(event shared.Event) {}
