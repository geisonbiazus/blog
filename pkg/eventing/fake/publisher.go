package fake

import "github.com/geisonbiazus/blog/pkg/eventing/entities"

type Publisher struct {
	Events []entities.Event
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Publish(event entities.Event) error {
	p.Events = append(p.Events, event)
	return nil
}

func (p *Publisher) LastEvent() entities.Event {
	if len(p.Events) == 0 {
		return entities.Event{}
	}

	return p.Events[len(p.Events)-1]
}

func (p *Publisher) Clear() {
	p.Events = []entities.Event{}
}
