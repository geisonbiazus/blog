package fake

import "context"

type TransactionManager struct{}

func NewFakeTransactionManager() *TransactionManager {
	return &TransactionManager{}
}

func (t *TransactionManager) Transaction(ctx context.Context, callback func(ctx context.Context) error) error {
	return callback(ctx)
}
