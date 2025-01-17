package ports

import (
	"context"
	"time"

	"github.com/geisonbiazus/blog/internal/auth/entities"
)

type OAuth2Provider interface {
	AuthURL(state string) string
	AuthenticatedUser(ctx context.Context, code string) (entities.ProviderUser, error)
}

type StateRepo interface {
	AddState(state string) error
	Exists(state string) (bool, error)
	Remove(state string) error
}

type UserRepo interface {
	CreateUser(ctx context.Context, user entities.User) error
	UpdateUser(ctx context.Context, user entities.User) error
	FindUserByProviderUserID(ctx context.Context, providerUserID string) (entities.User, error)
}

type TokenEncoder interface {
	Encode(value string, expiresIn time.Duration) (string, error)
}
