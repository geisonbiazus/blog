package dbrepo_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/geisonbiazus/blog/pkg/assert"
	"github.com/geisonbiazus/blog/pkg/dbrepo"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func TestTransactionManager(t *testing.T) {
	t.Run("Transaction", func(t *testing.T) {
		t.Run("It starts a transaction and puts it the callback context", func(t *testing.T) {
			db := dbrepo.ConnectoToTestDB()
			defer db.Close()

			manager := dbrepo.NewTransactionManager(db)
			var ctx context.Context
			manager.Transaction(context.Background(), func(c context.Context) error {
				ctx = c
				return nil
			})
			assert.NotNil(t, ctx)

			tx := dbrepo.TxFromContext(ctx)
			assert.NotNil(t, tx)
		})

		t.Run("It commits the transaction", func(t *testing.T) {
			db := dbrepo.ConnectoToTestDB()
			defer db.Close()

			createTestTable(db)

			manager := dbrepo.NewTransactionManager(db)
			manager.Transaction(context.Background(), func(ctx context.Context) error {
				tx := dbrepo.TxFromContext(ctx)
				insertValue(tx)
				return nil
			})

			assert.Equal(t, 1, countValues(db))

			dropTestTable(db)
		})

		t.Run("It rollbacks the transaction if the callback returns an error", func(t *testing.T) {
			db := dbrepo.ConnectoToTestDB()
			defer db.Close()

			createTestTable(db)

			manager := dbrepo.NewTransactionManager(db)
			manager.Transaction(context.Background(), func(ctx context.Context) error {
				tx := dbrepo.TxFromContext(ctx)
				insertValue(tx)
				return errors.New("some error")
			})

			assert.Equal(t, 0, countValues(db))

			dropTestTable(db)
		})

		t.Run("It rollbacks the transaction the callack panics", func(t *testing.T) {
			db := dbrepo.ConnectoToTestDB()
			defer db.Close()

			createTestTable(db)

			manager := dbrepo.NewTransactionManager(db)
			manager.Transaction(context.Background(), func(ctx context.Context) error {
				tx := dbrepo.TxFromContext(ctx)
				insertValue(tx)

				panic("error")
			})

			assert.Equal(t, 0, countValues(db))

			dropTestTable(db)
		})
	})
}

func createTestTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS test_tx_manager (col NUMERIC)")
	if err != nil {
		panic(err)
	}
}

func dropTestTable(db *sql.DB) {
	_, err := db.Exec("DROP TABLE test_tx_manager")
	if err != nil {
		panic(err)
	}
}

func insertValue(tx *sql.Tx) {
	_, err := tx.Exec("INSERT INTO test_tx_manager (col) VALUES (100)")
	if err != nil {
		panic(err)
	}
}

func countValues(db *sql.DB) int {
	var count int
	row := db.QueryRow("SELECT count(1) FROM test_tx_manager")
	err := row.Scan(&count)
	if err != nil {
		panic(err)
	}
	return count
}
