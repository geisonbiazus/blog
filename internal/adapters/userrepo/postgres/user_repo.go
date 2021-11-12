package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/geisonbiazus/blog/internal/core/auth"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Querier interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type TxKeyType string

var TxKey = TxKeyType("tx")

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) querier(ctx context.Context) Querier {
	tx := ctx.Value(TxKey)

	if tx == nil {
		return r.db
	}

	if tx, ok := tx.(*sql.Tx); ok {
		return tx
	}

	return r.db
}

func (r *UserRepo) CreateUser(ctx context.Context, user auth.User) error {
	q := r.querier(ctx)

	result, err := q.ExecContext(ctx,
		"INSERT INTO users (id, name, email, provider_user_id, avatar_url) VALUES ($1, $2, $3, $4, $5)",
		user.ID, user.Name, user.Email, user.ProviderUserID, user.AvatarURL,
	)

	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows != 1 {
		return errors.New("error inserting user")
	}

	return nil
}

func (r *UserRepo) FindById(ctx context.Context, id string) (auth.User, error) {
	q := r.querier(ctx)

	row := q.QueryRowContext(ctx, "SELECT id, name, email, provider_user_id, avatar_url FROM users where id = $1", id)

	user := auth.User{}

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.ProviderUserID, &user.AvatarURL)

	return user, err
}
