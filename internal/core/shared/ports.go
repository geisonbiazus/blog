package shared

import (
	"context"
	"time"
)

type TransactionManager interface {
	Transaction(ctx context.Context, callback func(ctx context.Context) error) error
}

type Cache[T any] interface {
	With(key string, expiresAt time.Time, resolve func() T) T
}
