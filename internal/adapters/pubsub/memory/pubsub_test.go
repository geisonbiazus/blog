package memory_test

import (
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/adapters/pubsub/memory"
	"github.com/geisonbiazus/blog/internal/core/shared"
	"github.com/stretchr/testify/suite"
)

type PubSubSuite struct {
	suite.Suite
	pubsub *memory.PubSub
}

func (s *PubSubSuite) SetupSubTest() {
	s.pubsub = memory.NewPubSub()
}

func (s *PubSubSuite) TestPublishAndSubscribe() {
	s.Run("It publishes and receives one event", func() {
		event := newTestEvent()

		subscriberChan := s.pubsub.Subscribe(event.Type)

		s.pubsub.Publish(event)

		s.Equal(event, <-subscriberChan)
	})

	s.Run("It does nothing when there is no subscriber", func() {
		event := newTestEvent()
		s.Nil(s.pubsub.Publish(event))
	})

	s.Run("It receives the published events in order", func() {
		event1 := newTestEvent()

		event2 := newTestEvent()
		event2.OccurredOn = time.Now().Add(-2 * time.Second)
		event2.Payload["key"] = "value2"

		subscriberChan := s.pubsub.Subscribe(event1.Type)

		s.pubsub.Publish(event1)
		s.pubsub.Publish(event2)

		s.Equal(event1, <-subscriberChan)
		s.Equal(event2, <-subscriberChan)
	})

	s.Run("It allows to publish new events after receiveiving the previous ones", func() {
		event1 := newTestEvent()

		event2 := newTestEvent()
		event2.OccurredOn = time.Now().Add(-2 * time.Second)
		event2.Payload["key"] = "value2"

		subscriberChan := s.pubsub.Subscribe(event1.Type)

		s.pubsub.Publish(event1)

		s.Equal(event1, <-subscriberChan)

		s.pubsub.Publish(event2)

		s.Equal(event2, <-subscriberChan)
	})

	s.Run("It subscribes to different event types", func() {
		event1 := newTestEvent()
		event1.Type = "EventType1"

		event2 := newTestEvent()
		event2.Type = "EventType2"

		subscriberChan1 := s.pubsub.Subscribe(event1.Type)
		subscriberChan2 := s.pubsub.Subscribe(event2.Type)

		s.pubsub.Publish(event1)
		s.pubsub.Publish(event2)

		s.Equal(event2, <-subscriberChan2)
		s.Equal(event1, <-subscriberChan1)
	})

	s.Run("It allows to have multiple subscribers for the same event type", func() {
		event := newTestEvent()

		subscriberChan1 := s.pubsub.Subscribe(event.Type)
		subscriberChan2 := s.pubsub.Subscribe(event.Type)

		s.pubsub.Publish(event)

		s.Equal(event, <-subscriberChan1)
		s.Equal(event, <-subscriberChan2)
	})
}

func newTestEvent() shared.Event {
	return shared.Event{
		Type:       "TestEvent",
		OccurredOn: time.Now(),
		Payload: map[string]interface{}{
			"key": "value1",
		},
	}
}

func TestPubSubSuite(t *testing.T) {
	suite.Run(t, &PubSubSuite{})
}
