package transaction

import (
	"context"
	"database/sql"

	"github.com/geisonbiazus/blog/pkg/transaction/fake"
	"github.com/geisonbiazus/blog/pkg/transaction/postgres"
)

func NewFakeManager() *fake.TransactionManager {
	return fake.NewTransactionManager()
}

func NewPostgresManager(db *sql.DB) *postgres.TransactionManager {
	return postgres.NewTransactionManager(db)
}

type Manager interface {
	Transaction(ctx context.Context, callback func(ctx context.Context) error) error
}
