package subscriptions

import (
	"context"

	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/internal/core/discussion"
)

type Subscriptions struct {
	subscriber Subscriber
	usecases   *UseCases
}

func New(subscriber Subscriber, usecases *UseCases) *Subscriptions {
	return &Subscriptions{subscriber: subscriber, usecases: usecases}
}

func (s *Subscriptions) Start() {
	NewSaveAuthorSubscriber(s.usecases.SaveAuthor, s.subscriber).Start()

	go func() {
		for event := range s.subscriber.Subscribe(auth.UserUpdatedEvent) {
			s.usecases.SaveAuthor.Run(context.Background(), discussion.SaveAuthorInput{
				ID:        event.Payload["ID"].(string),
				Name:      event.Payload["Name"].(string),
				AvatarURL: event.Payload["AvatarURL"].(string),
			})
		}
	}()
}
