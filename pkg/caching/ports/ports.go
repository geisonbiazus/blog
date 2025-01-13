package ports

import "time"

type ResolveFn func() (interface{}, error)

type Cache interface {
	Do(key string, resolve ResolveFn, expiresIn time.Duration) (interface{}, error)
}

var NeverExpire time.Duration = 0
