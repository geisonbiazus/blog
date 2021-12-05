package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
)

type Connection interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type Base struct {
	db *sql.DB
}

func NewBase(db *sql.DB) *Base {
	return &Base{db: db}
}

func (r *Base) Conn(ctx context.Context) Connection {
	tx := TxFromContext(ctx)

	if tx != nil {
		return tx
	}

	return r.db
}

func (r *Base) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	conn := r.Conn(ctx)

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
