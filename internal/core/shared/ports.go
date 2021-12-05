package shared

import "context"

type TransactionManager interface {
	Transaction(ctx context.Context, callback func(ctx context.Context) error) error
}
