package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/geisonbiazus/blog/pkg/dbrepo"
)

type TransactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (t *TransactionManager) Transaction(ctx context.Context, callback func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error on Transaction when beginning transaction: %w", err)
	}

	ctx = dbrepo.AddTxToContext(ctx, tx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = callback(ctx)

	if err != nil {
		e := tx.Rollback()
		if e != nil {
			return fmt.Errorf("error on Transaction when rolling back transaction: %w", e)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error on Transaction when commiting transaction: %w", err)
	}
	return nil
}
