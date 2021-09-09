package auth

import "errors"

var ErrInvalidState = errors.New("invalid state error")

type ProviderUser struct {
	ID        string
	Email     string
	Name      string
	AvatarURL string
}

type User struct {
	ID             string
	ProviderUserID string
	Email          string
	Name           string
	AvatarURL      string
}

var ErrUserNotFound = errors.New("user not found")
