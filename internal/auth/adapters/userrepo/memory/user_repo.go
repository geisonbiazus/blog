package memory

import (
	"context"

	"github.com/geisonbiazus/blog/internal/auth/entities"
)

type UserRepo struct {
	users []entities.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{users: []entities.User{}}
}

func (r *UserRepo) CreateUser(ctx context.Context, user entities.User) error {
	r.users = append(r.users, user)
	return nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, user entities.User) error {
	for i, existingUser := range r.users {
		if existingUser.ID == user.ID {
			r.users[i] = user
			return nil
		}
	}
	return entities.ErrUserNotFound
}

func (r *UserRepo) FindUserByID(ctx context.Context, id string) (entities.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}

	return entities.User{}, entities.ErrUserNotFound
}

func (r *UserRepo) FindUserByEmail(ctx context.Context, email string) (entities.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return entities.User{}, entities.ErrUserNotFound
}

func (r *UserRepo) FindUserByProviderUserID(ctx context.Context, providerUserID string) (entities.User, error) {
	for _, user := range r.users {
		if user.ProviderUserID == providerUserID {
			return user, nil
		}
	}

	return entities.User{}, entities.ErrUserNotFound
}
