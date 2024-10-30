package shared

import (
	"context"
	"time"
)

type TransactionManager interface {
	Transaction(ctx context.Context, callback func(ctx context.Context) error) error
}

type ResolveFn func() (interface{}, error)

type Cache interface {
	Do(key string, resolve ResolveFn, expiresIn time.Duration) (interface{}, error)
}

var NeverExpire time.Duration = 0

type Publisher interface {
	Publish(event Event) error
}

type IDGenerator interface {
	Generate() string
}
