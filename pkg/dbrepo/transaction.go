package dbrepo

import (
	"context"
	"database/sql"
)

type TxKeyType string

var TxKey = TxKeyType("tx")

func AddTxToContext(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

func TxFromContext(ctx context.Context) *sql.Tx {
	val := ctx.Value(TxKey)

	if val == nil {
		return nil
	}

	tx, ok := val.(*sql.Tx)
	if !ok {
		return nil
	}

	return tx
}
