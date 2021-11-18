package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/geisonbiazus/blog/internal/core/auth"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Connection interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type TxKeyType string

var TxKey = TxKeyType("tx")

type DBRepo struct {
	db *sql.DB
}

func (r *DBRepo) conn(ctx context.Context) Connection {
	tx := ctx.Value(TxKey)

	if tx == nil {
		return r.db
	}

	if tx, ok := tx.(*sql.Tx); ok {
		return tx
	}

	return r.db
}

func (r *DBRepo) exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	conn := r.conn(ctx)

	result, err := conn.ExecContext(ctx, query, args...)

	if err != nil {
		return 0, fmt.Errorf("error on exec when executing query: %w", err)
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("error on exec when getting affected rows: %w", err)
	}

	return rows, nil
}

type UserRepo struct {
	DBRepo
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DBRepo: DBRepo{db: db}}
}

func (r *UserRepo) CreateUser(ctx context.Context, user auth.User) error {
	rows, err := r.exec(ctx,
		"INSERT INTO users (id, name, email, provider_user_id, avatar_url) VALUES ($1, $2, $3, $4, $5)",
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
	rows, err := r.exec(ctx, `
		UPDATE users set 
			name = $1, 
			email = $2, 
			provider_user_id = $3, 
			avatar_url = $4
		WHERE id = $5`,
		user.Name, user.Email, user.ProviderUserID, user.AvatarURL, user.ID,
	)

	if err != nil {
		return fmt.Errorf("error on Update when executing query: %w", err)
	}

	if rows != 1 {
		return auth.ErrUserNotFound
	}

	return nil
}

func (r *UserRepo) FindUserByID(ctx context.Context, id string) (auth.User, error) {
	conn := r.conn(ctx)

	row := conn.QueryRowContext(ctx, "SELECT id, name, email, provider_user_id, avatar_url FROM users where id = $1", id)

	user := auth.User{}

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.ProviderUserID, &user.AvatarURL)

	return user, err
}
