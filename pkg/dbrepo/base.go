package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var ErrNoRowsUpdated = errors.New("error on Update, no affected rows")

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

func (r *Base) Insert(
	ctx context.Context, tableName string, values map[string]interface{},
) error {
	columns := []string{}
	placeholders := []string{}
	args := []interface{}{}

	for key, value := range values {
		columns = append(columns, key)
		args = append(args, value)
	}

	for i := range columns {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
	}

	statement := fmt.Sprintf(`
		INSERT INTO %s 
			(%s) 
		VALUES 
			(%s)`,
		tableName,
		strings.Join(columns, ","),
		strings.Join(placeholders, ","),
	)

	rows, err := r.Exec(ctx, statement, args...)

	if err != nil {
		return fmt.Errorf("error on Insert when executing query: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("error on Insert, no affected rows")
	}

	return nil
}

func (r *Base) Update(
	ctx context.Context, tableName string, id string, values map[string]interface{},
) error {
	columns := []string{}
	args := []interface{}{id}

	lastArg := 1

	for key, value := range values {
		columns = append(columns, fmt.Sprintf("%s = $%d", key, lastArg+1))
		args = append(args, value)
		lastArg++
	}

	statement := fmt.Sprintf(`
		UPDATE %s 
		SET %s 
		WHERE id = $1`,
		tableName,
		strings.Join(columns, ","),
	)

	rows, err := r.Exec(ctx, statement, args...)

	if err != nil {
		return fmt.Errorf("error on Update when executing query: %w", err)
	}

	if rows != 1 {
		return ErrNoRowsUpdated
	}

	return nil
}
