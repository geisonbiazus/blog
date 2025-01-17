package shared

import (
	"database/sql"

	"github.com/geisonbiazus/blog/pkg/caching"
	"github.com/geisonbiazus/blog/pkg/env"
	"github.com/geisonbiazus/blog/pkg/eventing"
	"github.com/geisonbiazus/blog/pkg/eventing/memory"
	"github.com/geisonbiazus/blog/pkg/gen"
	"github.com/geisonbiazus/blog/pkg/transaction"
)

type Context struct {
	Env string

	BaseURL string

	PostgresURL     string
	PostgresTestURL string

	db                 *sql.DB
	transactionManager transaction.Manager
	pubsub             *memory.PubSub
	cache              caching.Cache
}

func NewContext() *Context {
	return &Context{
		Env: env.GetString("ENV", "development"),

		BaseURL: env.GetString("BASE_URL", "http://localhost:3000"),

		PostgresURL:     env.GetString("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/blog?sslmode=disable"),
		PostgresTestURL: env.GetString("POSTGRES_TEST_URL", "postgres://postgres:postgres@localhost:5433/blog_test?sslmode=disable"),
	}
}

func (c *Context) IDGenerator() gen.Generator {
	return gen.NewUUIDGenerator()
}

func (c *Context) TransactionManager() transaction.Manager {
	if c.transactionManager == nil {
		c.transactionManager = c.resolveTransactionManager()
	}
	return c.transactionManager
}

func (c *Context) resolveTransactionManager() transaction.Manager {
	tm := transaction.NewPostgresManager(c.DB())
	if c.IsTest() {
		tm.EnableTestMode()
	}
	return tm
}

func (c *Context) DB() *sql.DB {
	if c.db == nil {
		db, err := sql.Open("pgx", c.resolvePostgresURL())
		if err != nil {
			panic(err)
		}
		c.db = db
	}
	return c.db
}

func (c *Context) resolvePostgresURL() string {
	if c.IsTest() {
		return c.PostgresTestURL
	}
	return c.PostgresURL
}

func (c *Context) Publisher() eventing.Publisher {
	return c.PubSub()
}

func (c *Context) PubSub() *memory.PubSub {
	if c.pubsub == nil {
		c.pubsub = eventing.NewMemoryPubSub()
	}
	return c.pubsub
}

func (c *Context) Cache() caching.Cache {
	if c.cache == nil {
		c.cache = c.resolveCache()
	}
	return c.cache
}

func (c *Context) resolveCache() caching.Cache {
	if c.IsDevelopment() {
		return caching.NewNullCache()
	}
	return caching.NewMemoryCache()
}

// Helpers

func (c *Context) IsTest() bool {
	return c.Env == "test"
}

func (c *Context) IsDevelopment() bool {
	return c.Env == "development"
}
