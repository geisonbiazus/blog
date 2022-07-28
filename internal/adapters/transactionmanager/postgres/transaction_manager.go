package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/geisonbiazus/blog/pkg/dbrepo"
)

type TransactionManager struct {
	db       *sql.DB
	testMode bool
}

func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (t *TransactionManager) EnableTestMode() {
	t.testMode = true
}

func (t *TransactionManager) Transaction(ctx context.Context, callback func(ctx context.Context) error) error {
	if t.isAlreadyInTransaction(ctx) {
		return callback(ctx)
	}

	tx, err := t.beginTransaction(ctx)
	if err != nil {
		return err
	}

	defer t.rollbackOnPanic(tx)

	ctx = dbrepo.AddTxToContext(ctx, tx)

	err = callback(ctx)
	if err != nil {
		return t.rollbackTransaction(tx, err)
	}

	if t.testMode {
		return t.rollbackTransaction(tx, nil)
	}

	return t.commitTransaction(tx)
}

func (t *TransactionManager) isAlreadyInTransaction(ctx context.Context) bool {
	return dbrepo.TxFromContext(ctx) != nil
}

func (t *TransactionManager) beginTransaction(ctx context.Context) (*sql.Tx, error) {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return tx, fmt.Errorf("error on Transaction when beginning transaction: %w", err)
	}

	return tx, nil
}

func (t *TransactionManager) rollbackOnPanic(tx *sql.Tx) {
	if r := recover(); r != nil {
		err := t.rollbackTransaction(tx, nil)
		if err != nil {
			panic(err)
		}
		panic(r)
	}
}

func (t *TransactionManager) rollbackTransaction(tx *sql.Tx, err error) error {
	e := tx.Rollback()
	if e != nil {
		return fmt.Errorf("error on Transaction when rolling back transaction: %w", e)
	}
	return err
}

func (t *TransactionManager) commitTransaction(tx *sql.Tx) error {
	err := tx.Commit()
	if err != nil {
		return fmt.Errorf("error on Transaction when commiting transaction: %w", err)
	}
	return nil
}
