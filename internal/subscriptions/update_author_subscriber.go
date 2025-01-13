package subscriptions

import (
	"context"

	"github.com/geisonbiazus/blog/internal/auth"
	"github.com/geisonbiazus/blog/internal/discussion"
	"github.com/geisonbiazus/blog/pkg/eventing"
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
	s.BaseSubscriber.Start(func(event eventing.Event) error {
		_, err := s.usecase.Run(context.Background(), s.inputFrom(event))
		return err
	})
}

func (s *UpdateAuthorSubscriber) inputFrom(event eventing.Event) discussion.SaveAuthorInput {
	return discussion.SaveAuthorInput{
		UserID:    event.Payload["ID"].(string),
		Name:      event.Payload["Name"].(string),
		AvatarURL: event.Payload["AvatarURL"].(string),
	}
}
