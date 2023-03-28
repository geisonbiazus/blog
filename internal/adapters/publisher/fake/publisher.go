package fake

import "github.com/geisonbiazus/blog/internal/core/shared"

type Publisher struct {
	Events []shared.Event
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Publish(event shared.Event) error {
	p.Events = append(p.Events, event)
	return nil
}

func (p *Publisher) LastEvent() shared.Event {
	if len(p.Events) == 0 {
		return shared.Event{}
	}

	return p.Events[len(p.Events)-1]
}

func (p *Publisher) Clear() {
	p.Events = []shared.Event{}
}
