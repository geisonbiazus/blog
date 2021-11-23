package dbrepo

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// Test connects to the database, starts a transaction and puts it
// in the context passing it to the given callback. This context is supposed
// to be used by all the repositories in the test.
// After the callback is run, the transaction is rolled back returning
// the database to the initial state.
func Test(cb func(ctx context.Context, db *sql.DB)) {
	ctx := context.Background()
	db := ConnectoToTestDB()
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	ctx = context.WithValue(ctx, TxKey, tx)
	cb(ctx, db)
	tx.Rollback()
}

func ConnectoToTestDB() *sql.DB {
	db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5433/blog_test?sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db
}
