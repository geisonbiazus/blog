package subscriptions

import (
	"context"

	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/geisonbiazus/blog/internal/core/shared"
)

type SaveAuthorSubscriber struct {
	*BaseSubscriber
	usecase SaveAuthorUseCase
}

func NewSaveAuthorSubscriber(usecase SaveAuthorUseCase, subscriber Subscriber) *SaveAuthorSubscriber {
	return &SaveAuthorSubscriber{
		BaseSubscriber: NewBaseSubscriber(subscriber, auth.UserCreatedEvent),
		usecase:        usecase,
	}
}

func (s *SaveAuthorSubscriber) Start() {
	s.BaseSubscriber.Start(func(event shared.Event) error {
		_, err := s.usecase.Run(context.Background(), s.inputFrom(event))
		return err
	})
}

func (s *SaveAuthorSubscriber) inputFrom(event shared.Event) discussion.SaveAuthorInput {
	return discussion.SaveAuthorInput{
		UserID:    event.Payload["ID"].(string),
		Name:      event.Payload["Name"].(string),
		AvatarURL: event.Payload["AvatarURL"].(string),
	}
}
