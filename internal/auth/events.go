package auth

import (
	"time"

	"github.com/geisonbiazus/blog/pkg/eventing"
)

const (
	UserCreatedEvent = "UserCreated"
	UserUpdatedEvent = "UserUpdated"
)

func NewUserCreatedEvent(user User) eventing.Event {
	return eventing.Event{
		Type:       UserCreatedEvent,
		OccurredOn: time.Now(),
		Payload: map[string]interface{}{
			"ID":        user.ID,
			"Email":     user.Email,
			"Name":      user.Name,
			"AvatarURL": user.AvatarURL,
		},
	}
}

func NewUserUpdatedEvent(user User) eventing.Event {
	return eventing.Event{
		Type:       UserUpdatedEvent,
		OccurredOn: time.Now(),
		Payload: map[string]interface{}{
			"ID":        user.ID,
			"Email":     user.Email,
			"Name":      user.Name,
			"AvatarURL": user.AvatarURL,
		},
	}
}
