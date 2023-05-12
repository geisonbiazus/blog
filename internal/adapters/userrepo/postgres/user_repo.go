package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/pkg/dbrepo"
)

type UserRepo struct {
	*dbrepo.Base
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{Base: dbrepo.NewBase(db)}
}

func (r *UserRepo) CreateUser(ctx context.Context, user auth.User) error {
	err := r.Insert(ctx, "auth_users", map[string]interface{}{
		"id":               user.ID,
		"name":             user.Name,
		"email":            user.Email,
		"provider_user_id": user.ProviderUserID,
		"avatar_url":       user.AvatarURL,
	})

	if err != nil {
		return fmt.Errorf("error on CreateUser: %w", err)
	}

	return nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, user auth.User) error {
	err := r.Update(ctx, "auth_users", user.ID, map[string]interface{}{
		"name":             user.Name,
		"email":            user.Email,
		"provider_user_id": user.ProviderUserID,
		"avatar_url":       user.AvatarURL,
	})

	if err != nil {
		if errors.Is(err, dbrepo.ErrNoRowsUpdated) {
			return auth.ErrUserNotFound
		}

		return fmt.Errorf("error on UpdateAuthor: %w", err)
	}

	return nil
}

func (r *UserRepo) FindUserByID(ctx context.Context, id string) (auth.User, error) {
	return r.findUserBy(ctx, "id", id)
}

func (r *UserRepo) FindUserByProviderUserID(ctx context.Context, providerUserID string) (auth.User, error) {
	return r.findUserBy(ctx, "provider_user_id", providerUserID)
}

func (r *UserRepo) findUserBy(ctx context.Context, field string, value interface{}) (auth.User, error) {
	conn := r.Conn(ctx)

	row := conn.QueryRowContext(ctx, `
		SELECT 
			id, name, email, provider_user_id, avatar_url 
		FROM auth_users 
		WHERE `+field+` = $1`,
		value,
	)

	user := auth.User{}

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.ProviderUserID, &user.AvatarURL)

	if errors.Is(err, sql.ErrNoRows) {
		return auth.User{}, auth.ErrUserNotFound
	}

	if err != nil {
		return auth.User{}, fmt.Errorf("error on findUserBy when executing query: %w", err)
	}

	return user, err
}
