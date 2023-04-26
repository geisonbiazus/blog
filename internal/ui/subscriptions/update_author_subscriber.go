package subscriptions

import (
	"context"

	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/geisonbiazus/blog/internal/core/shared"
)

type UpdateAuthorSubscriber struct {
	*BaseSubscriber
	usecase SaveAuthorUseCase
}

func NewUpdateAuthorSubscriber(usecase SaveAuthorUseCase, subscriber Subscriber) *UpdateAuthorSubscriber {
	return &UpdateAuthorSubscriber{
		BaseSubscriber: NewBaseSubscriber(subscriber, auth.UserUpdatedEvent),
		usecase:        usecase,
	}
}

func (s *UpdateAuthorSubscriber) Start() {
	s.BaseSubscriber.Start(func(event shared.Event) error {
		_, err := s.usecase.Run(context.Background(), s.inputFrom(event))
		return err
	})
}

func (s *UpdateAuthorSubscriber) inputFrom(event shared.Event) discussion.SaveAuthorInput {
	return discussion.SaveAuthorInput{
		UserID:    event.Payload["ID"].(string),
		Name:      event.Payload["Name"].(string),
		AvatarURL: event.Payload["AvatarURL"].(string),
	}
}
