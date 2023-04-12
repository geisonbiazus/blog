package subscriptions_test

import (
	"errors"
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/geisonbiazus/blog/internal/core/shared"
	"github.com/geisonbiazus/blog/internal/ui/subscriptions"
	"github.com/stretchr/testify/suite"
)

type UpdateAuthorSubscriberSuite struct {
	suite.Suite
	updateAuthorSubscriber *subscriptions.UpdateAuthorSubscriber
	usecase                *SaveAuthorUseCaseSpy
	subscriber             *SubscriberSpy
	event                  shared.Event
}

func (s *UpdateAuthorSubscriberSuite) SetupSubTest() {
	s.usecase = NewSaveAuthorUseCaseSpy()
	s.subscriber = NewSubscriberSpy()
	s.updateAuthorSubscriber = subscriptions.NewUpdateAuthorSubscriber(s.usecase, s.subscriber)
	s.event = shared.Event{
		Type:       auth.UserUpdatedEvent,
		OccurredOn: time.Now(),
		Payload: map[string]interface{}{
			"ID":        "ID",
			"Email":     "user@example.com",
			"Name":      "Name",
			"AvatarURL": "http://example.com/avatar.png",
		},
	}
}

func (s *UpdateAuthorSubscriberSuite) TestStart() {
	s.Run("It executes the usecase when UserCreated event is publised", func() {
		s.updateAuthorSubscriber.Start()
		s.subscriber.Publish(s.event)

		s.True(<-s.usecase.Ran)
		s.Equal(discussion.SaveAuthorInput{
			ID:        "ID",
			Name:      "Name",
			AvatarURL: "http://example.com/avatar.png",
		}, s.usecase.ReceivedInput)
	})

	s.Run("It notifies success execution", func() {
		s.updateAuthorSubscriber.Start()
		s.subscriber.Publish(s.event)

		s.True(<-s.usecase.Ran)
		s.True(<-s.subscriber.Notified)
		s.Equal(s.event, s.subscriber.NotifySuccessReceivedEvent)
	})

	s.Run("It notifies the subscriber in case of error", func() {
		s.usecase.ReturnError = errors.New("error")
		s.updateAuthorSubscriber.Start()
		s.subscriber.Publish(s.event)

		s.True(<-s.usecase.Ran)
		s.True(<-s.subscriber.Notified)
		s.Equal(s.event, s.subscriber.NotifyErrorReceivedEvent)
		s.Equal(s.usecase.ReturnError, s.subscriber.NotifyErrorReceivedError)
	})
}

func TestUpdateAuthorSubscriberSuite(t *testing.T) {
	suite.Run(t, new(UpdateAuthorSubscriberSuite))
}
