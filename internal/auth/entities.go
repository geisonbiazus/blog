package auth

import "errors"

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

var ErrInvalidState = errors.New("invalid state error")
var ErrUserNotFound = errors.New("user not found")
var ErrTokenExpired = errors.New("token expired")
