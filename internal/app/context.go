package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/geisonbiazus/blog/internal/app/shared"
	"github.com/geisonbiazus/blog/internal/auth"
	"github.com/geisonbiazus/blog/internal/blog"
	"github.com/geisonbiazus/blog/internal/discussion"
	"github.com/geisonbiazus/blog/internal/discussion/adapters/commentrepo"
	"github.com/geisonbiazus/blog/internal/subscriptions"
	"github.com/geisonbiazus/blog/internal/web"
	webports "github.com/geisonbiazus/blog/internal/web/ports"
	"github.com/geisonbiazus/blog/pkg/env"
	"github.com/geisonbiazus/blog/pkg/eventing"
	"github.com/geisonbiazus/blog/pkg/eventing/memory"
	"github.com/geisonbiazus/blog/pkg/gen"
	"github.com/geisonbiazus/blog/pkg/migration"
	"github.com/geisonbiazus/blog/pkg/transaction"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Context struct {
	Env string

	Port           int
	TemplatePath   string
	StaticPath     string
	PostPath       string
	MigrationsPath string
	BaseURL        string

	PostgresURL     string
	PostgresTestURL string

	db                 *sql.DB
	transactionManager transaction.Manager
	pubsub             *memory.PubSub
	commentRepo        discussion.CommentRepo

	sharedContext *shared.Context
	blog          *blog.Context
	auth          *auth.Context
	discussion    *discussion.Context
}

func NewContext() *Context {
	return &Context{
		Env: env.GetString("ENV", "development"),

		Port:           env.GetInt("PORT", 3000),
		TemplatePath:   env.GetString("TEMPLATE_PATH", filepath.Join("web", "template")),
		StaticPath:     env.GetString("STATIC_PATH", filepath.Join("web", "static")),
		MigrationsPath: env.GetString("MIGRATIONS_PATH", "file://"+filepath.Join("db", "migrations")),
		BaseURL:        env.GetString("BASE_URL", "http://localhost:3000"),

		PostgresURL:     env.GetString("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/blog?sslmode=disable"),
		PostgresTestURL: env.GetString("POSTGRES_TEST_URL", "postgres://postgres:postgres@localhost:5433/blog_test?sslmode=disable"),
	}
}

func (c *Context) SharedContext() *shared.Context {
	if c.sharedContext == nil {
		c.sharedContext = shared.NewContext()
	}
	return c.sharedContext
}

// Components
func (c *Context) Blog() *blog.Context {
	if c.blog == nil {
		c.blog = blog.NewContext(c.SharedContext())
	}
	return c.blog
}

func (c *Context) Auth() *auth.Context {
	if c.auth == nil {
		c.auth = auth.NewContext(c.SharedContext())
	}
	return c.auth
}

func (c *Context) Discussion() *discussion.Context {
	if c.discussion == nil {
		c.discussion = discussion.NewContext(c.SharedContext())
	}
	return c.discussion
}

// UI

func (c *Context) WebServer() *web.Server {
	return web.NewServer(c.Port, c.Router(), c.Logger())
}

func (c *Context) Router() http.Handler {
	return web.NewRouter(c.TemplatePath, c.StaticPath, c.UseCases(), c.BaseURL)
}

func (c *Context) Subscriptions() *subscriptions.Subscriptions {
	return subscriptions.New(c.PubSub(), c.SubscriptionUseCases())
}

// Use cases

func (c *Context) UseCases() *webports.UseCases {
	return &webports.UseCases{
		ViewPost:      c.Blog().ViewPostUseCase(),
		ListPosts:     c.Blog().ListPostsUseCase(),
		RequestOAuth2: c.Auth().RequestOAuth2UseCase(),
		ConfirmOAuth2: c.Auth().ConfirmOAuth2UseCase(),
		ListComments:  c.Discussion().ListCommentsUseCase(),
	}
}

func (c *Context) SubscriptionUseCases() *subscriptions.UseCases {
	return &subscriptions.UseCases{
		SaveAuthor: c.Discussion().SaveAuthorUseCase(),
	}
}

// Adapters

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
	if c.isTest() {
		return c.PostgresTestURL
	}
	return c.PostgresURL
}

func (c *Context) Migration() *migration.Migration {
	return migration.New(c.DB(), c.MigrationsPath)
}

func (c *Context) PubSub() *memory.PubSub {
	if c.pubsub == nil {
		c.pubsub = eventing.NewMemoryPubSub()
	}
	return c.pubsub
}

func (c *Context) IDGenerator() gen.Generator {
	return gen.NewUUIDGenerator()
}

func (c *Context) CommentRepo() discussion.CommentRepo {
	if c.commentRepo == nil {
		c.commentRepo = commentrepo.NewPostgresCommentRepo(c.DB())
	}
	return c.commentRepo
}

func (c *Context) Logger() *log.Logger {
	return log.New(os.Stdout, "web: ", log.Ldate|log.Ltime|log.LUTC)
}

// Helpers

func (c *Context) isTest() bool {
	return c.Env == "test"
}
