package subscriptions_test

import (
	"context"

	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/geisonbiazus/blog/internal/core/shared"
)

type SaveAuthorUseCaseSpy struct {
	Ran             chan bool
	ReceivedContext context.Context
	ReceivedInput   discussion.SaveAuthorInput
	ReturnAuthor    *discussion.Author
	ReturnError     error
}

func NewSaveAuthorUseCaseSpy() *SaveAuthorUseCaseSpy {
	return &SaveAuthorUseCaseSpy{
		Ran: make(chan bool),
	}
}

func (s *SaveAuthorUseCaseSpy) Run(ctx context.Context, input discussion.SaveAuthorInput) (*discussion.Author, error) {
	s.ReceivedContext = ctx
	s.ReceivedInput = input
	s.Ran <- true
	return s.ReturnAuthor, s.ReturnError
}

type SubscriberSpy struct {
	channel                    chan shared.Event
	NotifySuccessReceivedEvent shared.Event
	NotifyErrorReceivedEvent   shared.Event
	NotifyErrorReceivedError   error
	Notified                   chan bool
}

func NewSubscriberSpy() *SubscriberSpy {
	return &SubscriberSpy{
		Notified: make(chan bool),
	}
}

func (f *SubscriberSpy) Publish(event shared.Event) {
	f.channel <- event
}

func (f *SubscriberSpy) Subscribe(eventType string) chan shared.Event {
	f.channel = make(chan shared.Event)
	return f.channel
}

func (f *SubscriberSpy) NotifyError(event shared.Event, err error) {
	f.NotifyErrorReceivedEvent = event
	f.NotifyErrorReceivedError = err
	f.Notified <- true
}

func (f *SubscriberSpy) NotifySuccess(event shared.Event) {
	f.NotifySuccessReceivedEvent = event
	f.Notified <- true
}
