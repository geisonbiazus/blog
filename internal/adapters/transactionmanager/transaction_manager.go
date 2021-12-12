package transactionmanager

import (
	"database/sql"

	"github.com/geisonbiazus/blog/internal/adapters/transactionmanager/fake"
	"github.com/geisonbiazus/blog/internal/adapters/transactionmanager/postgres"
)

func NewFakeTransactionManager() *fake.TransactionManager {
	return fake.NewTransactionManager()
}

func NewPostgresTransactionManager(db *sql.DB) *postgres.TransactionManager {
	return postgres.NewTransactionManager(db)
}
