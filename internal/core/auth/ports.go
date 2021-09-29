package auth

import (
	"context"
	"time"
)

type OAuth2Provider interface {
	AuthURL(state string) string
	AuthenticatedUser(ctx context.Context, code string) (ProviderUser, error)
}

type IDGenerator interface {
	Generate() string
}

type StateRepo interface {
	AddState(state string) error
	Exists(state string) (bool, error)
	Remove(state string) error
}

type UserRepo interface {
	CreateUser(user User) error
	UpdateUser(user User) error
	FindUserByProviderUserID(providerUserID string) (User, error)
}

type TokenManager interface {
	Encode(userID string, expiresIn time.Duration) (string, error)
}
