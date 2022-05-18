package shared

import (
	"context"
	"time"
)

type TransactionManager interface {
	Transaction(ctx context.Context, callback func(ctx context.Context) error) error
}

type Cache[T any] interface {
	Do(key string, resolve func() (T, error), expiresAt time.Time) (T, error)
}

var NeverExpire time.Duration = 0
