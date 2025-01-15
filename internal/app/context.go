package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/geisonbiazus/blog/internal/auth"
	"github.com/geisonbiazus/blog/internal/auth/adapters/oauth2provider"
	"github.com/geisonbiazus/blog/internal/auth/adapters/staterepo"
	"github.com/geisonbiazus/blog/internal/auth/adapters/tokenencoder"
	"github.com/geisonbiazus/blog/internal/auth/adapters/userrepo"
	"github.com/geisonbiazus/blog/internal/blog"
	"github.com/geisonbiazus/blog/internal/discussion"
	"github.com/geisonbiazus/blog/internal/discussion/adapters/commentrepo"
	"github.com/geisonbiazus/blog/internal/subscriptions"
	"github.com/geisonbiazus/blog/internal/web"
	webports "github.com/geisonbiazus/blog/internal/web/ports"
	"github.com/geisonbiazus/blog/pkg/caching"
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

	GitHubClientID     string
	GitHubClientSecret string

	AuthTokenSecret string

	PostgresURL     string
	PostgresTestURL string

	db                 *sql.DB
	transactionManager transaction.Manager
	pubsub             *memory.PubSub
	cache              caching.Cache
	stateRepo          auth.StateRepo
	userRepo           auth.UserRepo
	commentRepo        discussion.CommentRepo

	blog *blog.Context
}

func NewContext() *Context {
	return &Context{
		Env: env.GetString("ENV", "development"),

		Port:           env.GetInt("PORT", 3000),
		TemplatePath:   env.GetString("TEMPLATE_PATH", filepath.Join("web", "template")),
		StaticPath:     env.GetString("STATIC_PATH", filepath.Join("web", "static")),
		MigrationsPath: env.GetString("MIGRATIONS_PATH", "file://"+filepath.Join("db", "migrations")),
		BaseURL:        env.GetString("BASE_URL", "http://localhost:3000"),

		GitHubClientID:     env.GetString("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: env.GetString("GITHUB_CLIENT_SECRET", ""),

		AuthTokenSecret: env.GetString("AUTH_TOKEN_SECRET", ""),

		PostgresURL:     env.GetString("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/blog?sslmode=disable"),
		PostgresTestURL: env.GetString("POSTGRES_TEST_URL", "postgres://postgres:postgres@localhost:5433/blog_test?sslmode=disable"),
	}
}

// Components
func (c *Context) Blog() *blog.Context {
	if c.blog == nil {
		c.blog = blog.NewContext(c.Cache)
	}
	return c.blog
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
		RequestOAuth2: c.RequestOAuth2UseCase(),
		ConfirmOAuth2: c.ConfirmOAuth2UseCase(),
		ListComments:  c.ListCommentsUseCase(),
	}
}

func (c *Context) SubscriptionUseCases() *subscriptions.UseCases {
	return &subscriptions.UseCases{
		SaveAuthor: c.SaveAuthorUseCase(),
	}
}

func (c *Context) RequestOAuth2UseCase() *auth.RequestOAuth2UseCase {
	return auth.NewRequestOAuth2UseCase(c.OAuth2Provider(), c.IDGenerator(), c.StateRepo())
}

func (c *Context) ConfirmOAuth2UseCase() *auth.ConfirmOAuth2UseCase {
	return auth.NewConfirmOAuth2UseCase(c.OAuth2Provider(), c.StateRepo(), c.UserRepo(), c.IDGenerator(), c.TokenEncoder(), c.TransactionManager(), c.PubSub())
}

func (c *Context) ListCommentsUseCase() *discussion.ListCommentsUseCase {
	return discussion.NewListCommentsUseCase(c.CommentRepo())
}

func (c *Context) SaveAuthorUseCase() *discussion.SaveAuthorUseCase {
	return discussion.NewSaveAuthorUseCase(c.CommentRepo(), c.TransactionManager(), c.IDGenerator())
}

// Adapters

func (c *Context) Cache() caching.Cache {
	if c.cache == nil {
		c.cache = c.resolveCache()
	}
	return c.cache
}

func (c *Context) resolveCache() caching.Cache {
	if c.isDevelopment() {
		return caching.NewNullCache()
	}
	return caching.NewMemoryCache()
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
	if c.isTest() {
		return c.PostgresTestURL
	}
	return c.PostgresURL
}

func (c *Context) Migration() *migration.Migration {
	return migration.New(c.DB(), c.MigrationsPath)
}

func (c *Context) TransactionManager() transaction.Manager {
	if c.transactionManager == nil {
		c.transactionManager = c.resolveTransactionManager()
	}
	return c.transactionManager
}

func (c *Context) resolveTransactionManager() transaction.Manager {
	tm := transaction.NewPostgresManager(c.DB())
	if c.isTest() {
		tm.EnableTestMode()
	}
	return tm
}

func (c *Context) PubSub() *memory.PubSub {
	if c.pubsub == nil {
		c.pubsub = eventing.NewMemoryPubSub()
	}
	return c.pubsub
}

func (c *Context) OAuth2Provider() auth.OAuth2Provider {
	if c.isTest() {
		return c.FakeOAuth2Provider()
	}
	return c.GithubOAuth2Provider()
}

func (c *Context) GithubOAuth2Provider() auth.OAuth2Provider {
	return oauth2provider.NewGithubProvider(c.GitHubClientID, c.GitHubClientSecret)
}

func (c *Context) FakeOAuth2Provider() auth.OAuth2Provider {
	return oauth2provider.NewFakeProvider(c.BaseURL)
}

func (c *Context) IDGenerator() gen.Generator {
	return gen.NewUUIDGenerator()
}

func (c *Context) StateRepo() auth.StateRepo {
	if c.stateRepo == nil {
		c.stateRepo = staterepo.NewMemoryStateRepo()
	}
	return c.stateRepo
}

func (c *Context) UserRepo() auth.UserRepo {
	if c.userRepo == nil {
		c.userRepo = userrepo.NewPostgresUserRepo(c.DB())
	}
	return c.userRepo
}

func (c *Context) CommentRepo() discussion.CommentRepo {
	if c.commentRepo == nil {
		c.commentRepo = commentrepo.NewPostgresCommentRepo(c.DB())
	}
	return c.commentRepo
}

func (c *Context) TokenEncoder() auth.TokenEncoder {
	return tokenencoder.NewJWTTokenEncoder(c.AuthTokenSecret)
}

func (c *Context) Logger() *log.Logger {
	return log.New(os.Stdout, "web: ", log.Ldate|log.Ltime|log.LUTC)
}

// Helpers

func (c *Context) isTest() bool {
	return c.Env == "test"
}

func (c *Context) isDevelopment() bool {
	return c.Env == "development"
}
