package memory

import "github.com/geisonbiazus/blog/pkg/eventing/entities"

type PubSub struct {
	SubscriberBuffer int
	subscribers      map[string][]chan entities.Event
}

func NewPubSub() *PubSub {
	return &PubSub{
		SubscriberBuffer: 10,
		subscribers:      map[string][]chan entities.Event{},
	}
}

func (p *PubSub) Publish(event entities.Event) error {
	for _, subscriber := range p.subscribers[event.Type] {
		subscriber <- event
	}

	return nil
}

func (p *PubSub) Subscribe(eventType string) chan entities.Event {
	channel := make(chan entities.Event, p.SubscriberBuffer)
	p.ensureSubscribersFor(eventType)
	p.subscribers[eventType] = append(p.subscribers[eventType], channel)

	return channel
}

func (p *PubSub) ensureSubscribersFor(eventType string) {
	if _, ok := p.subscribers[eventType]; !ok {
		p.subscribers[eventType] = []chan entities.Event{}
	}
}

func (p *PubSub) NotifyError(event entities.Event, err error) {}

func (p *PubSub) NotifySuccess(event entities.Event) {}
