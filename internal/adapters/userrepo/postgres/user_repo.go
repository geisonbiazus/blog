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

type UserRepo struct {
	DBRepo
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DBRepo: DBRepo{db: db}}
}

func (r *UserRepo) CreateUser(ctx context.Context, user auth.User) error {
	conn := r.conn(ctx)

	result, err := conn.ExecContext(ctx,
		"INSERT INTO users (id, name, email, provider_user_id, avatar_url) VALUES ($1, $2, $3, $4, $5)",
		user.ID, user.Name, user.Email, user.ProviderUserID, user.AvatarURL,
	)

	if err != nil {
		return fmt.Errorf("error on CreateUser when executing query: %w", err)
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("error on CreateUser when executing query: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("error on CreateUser when executing query")
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
