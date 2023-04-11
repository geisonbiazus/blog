package subscriptions

import (
	"context"

	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/geisonbiazus/blog/internal/core/shared"
)

type SaveAuthorSubscriber struct {
	usecase    SaveAuthorUseCase
	subscriber Subscriber
}

func NewSaveAuthorSubscriber(usecase SaveAuthorUseCase, subscriber Subscriber) *SaveAuthorSubscriber {
	return &SaveAuthorSubscriber{usecase: usecase, subscriber: subscriber}
}

func (s *SaveAuthorSubscriber) Start() {
	subscription := s.subscriber.Subscribe(auth.UserCreatedEvent)

	go func() {
		for event := range subscription {
			_, err := s.usecase.Run(context.Background(), s.inputFrom(event))

			if err != nil {
				s.subscriber.NotifyError(event, err)
			} else {
				s.subscriber.NotifySuccess(event)
			}
		}
	}()
}

func (s *SaveAuthorSubscriber) inputFrom(event shared.Event) discussion.SaveAuthorInput {
	return discussion.SaveAuthorInput{
		ID:        event.Payload["ID"].(string),
		Name:      event.Payload["Name"].(string),
		AvatarURL: event.Payload["AvatarURL"].(string),
	}
}
