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

type SaveAuthorSubscriberSuite struct {
	suite.Suite
	saveAuthorSubscriber *subscriptions.SaveAuthorSubscriber
	usecase              *SaveAuthorUseCaseSpy
	subscriber           *SubscriberSpy
	event                shared.Event
}

func (s *SaveAuthorSubscriberSuite) SetupSubTest() {
	s.usecase = NewSaveAuthorUseCaseSpy()
	s.subscriber = NewSubscriberSpy()
	s.saveAuthorSubscriber = subscriptions.NewSaveAuthorSubscriber(s.usecase, s.subscriber)
	s.event = shared.Event{
		Type:       auth.UserCreatedEvent,
		OccurredOn: time.Now(),
		Payload: map[string]interface{}{
			"ID":        "ID",
			"Email":     "user@example.com",
			"Name":      "Name",
			"AvatarURL": "http://example.com/avatar.png",
		},
	}
}

func (s *SaveAuthorSubscriberSuite) TestStart() {
	s.Run("It executes the usecase when UserCreated event is publised", func() {
		s.saveAuthorSubscriber.Start()
		s.subscriber.Publish(s.event)

		s.True(<-s.usecase.Ran)
		s.Equal(discussion.SaveAuthorInput{
			ID:        "ID",
			Name:      "Name",
			AvatarURL: "http://example.com/avatar.png",
		}, s.usecase.ReceivedInput)
	})

	s.Run("It notifies success execution", func() {
		s.saveAuthorSubscriber.Start()
		s.subscriber.Publish(s.event)

		s.True(<-s.usecase.Ran)
		s.True(<-s.subscriber.Notified)
		s.Equal(s.event, s.subscriber.NotifySuccessReceivedEvent)
	})

	s.Run("It notifies the subscriber in case of error", func() {
		s.usecase.ReturnError = errors.New("error")
		s.saveAuthorSubscriber.Start()
		s.subscriber.Publish(s.event)

		s.True(<-s.usecase.Ran)
		s.True(<-s.subscriber.Notified)
		s.Equal(s.event, s.subscriber.NotifyErrorReceivedEvent)
		s.Equal(s.usecase.ReturnError, s.subscriber.NotifyErrorReceivedError)
	})
}

func TestSaveAuthorSubscriberSuite(t *testing.T) {
	suite.Run(t, new(SaveAuthorSubscriberSuite))
}
