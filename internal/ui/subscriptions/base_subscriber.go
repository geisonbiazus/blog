package subscriptions

import "github.com/geisonbiazus/blog/internal/core/shared"

type BaseSubscriber struct {
	subscriber Subscriber
	eventType  string
}

func NewBaseSubscriber(subscriber Subscriber, eventType string) *BaseSubscriber {
	return &BaseSubscriber{subscriber: subscriber, eventType: eventType}
}

func (s *BaseSubscriber) Start(runUseCase func(event shared.Event) error) {
	subscription := s.subscriber.Subscribe(s.eventType)

	go func() {
		for event := range subscription {
			err := runUseCase(event)

			if err != nil {
				s.subscriber.NotifyError(event, err)
			} else {
				s.subscriber.NotifySuccess(event)
			}
		}
	}()
}
