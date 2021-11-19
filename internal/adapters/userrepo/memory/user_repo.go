package memory

import (
	"context"

	"github.com/geisonbiazus/blog/internal/core/auth"
)

type UserRepo struct {
	users []auth.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{users: []auth.User{}}
}

func (r *UserRepo) CreateUser(ctx context.Context, user auth.User) error {
	r.users = append(r.users, user)
	return nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, user auth.User) error {
	for i, existingUser := range r.users {
		if existingUser.ID == user.ID {
			r.users[i] = user
			return nil
		}
	}
	return auth.ErrUserNotFound
}

func (r *UserRepo) FindUserByID(ctx context.Context, id string) (auth.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}

	return auth.User{}, auth.ErrUserNotFound
}

func (r *UserRepo) FindUserByEmail(ctx context.Context, email string) (auth.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return auth.User{}, auth.ErrUserNotFound
}

func (r *UserRepo) FindUserByProviderUserID(ctx context.Context, providerUserID string) (auth.User, error) {
	for _, user := range r.users {
		if user.ProviderUserID == providerUserID {
			return user, nil
		}
	}

	return auth.User{}, auth.ErrUserNotFound
}
