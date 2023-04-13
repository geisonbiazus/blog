package auth

import (
	"time"

	"github.com/geisonbiazus/blog/internal/core/shared"
)

const (
	UserCreatedEvent = "UserCreated"
	UserUpdatedEvent = "UserUpdated"
)

func NewUserCreatedEvent(user User) shared.Event {
	return shared.Event{
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

func NewUserUpdatedEvent(user User) shared.Event {
	return shared.Event{
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
