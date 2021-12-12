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
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user User) error
	FindUserByProviderUserID(ctx context.Context, providerUserID string) (User, error)
}

type TokenEncoder interface {
	Encode(value string, expiresIn time.Duration) (string, error)
}
