package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/pkg/dbrepo"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type UserRepo struct {
	*dbrepo.Base
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{Base: dbrepo.NewBase(db)}
}

func (r *UserRepo) CreateUser(ctx context.Context, user auth.User) error {
	rows, err := r.Exec(ctx, `
		INSERT INTO users 
			(id, name, email, provider_user_id, avatar_url) 
		VALUES 
			($1, $2, $3, $4, $5)`,
		user.ID, user.Name, user.Email, user.ProviderUserID, user.AvatarURL,
	)

	if err != nil {
		return fmt.Errorf("error on CreateUser when executing query: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("error on CreateUser, no affected rows")
	}

	return nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, user auth.User) error {
	rows, err := r.Exec(ctx, `
		UPDATE users set 
			name = $1, 
			email = $2, 
			provider_user_id = $3, 
			avatar_url = $4
		WHERE id = $5`,
		user.Name, user.Email, user.ProviderUserID, user.AvatarURL, user.ID,
	)

	if err != nil {
		return fmt.Errorf("error on UpdateUser when executing query: %w", err)
	}

	if rows != 1 {
		return auth.ErrUserNotFound
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
		FROM users 
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
