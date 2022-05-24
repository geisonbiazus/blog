package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/geisonbiazus/blog/internal/adapters/cache"
	"github.com/geisonbiazus/blog/internal/adapters/idgenerator"
	"github.com/geisonbiazus/blog/internal/adapters/oauth2provider"
	"github.com/geisonbiazus/blog/internal/adapters/postrepo"
	"github.com/geisonbiazus/blog/internal/adapters/renderer"
	"github.com/geisonbiazus/blog/internal/adapters/staterepo"
	"github.com/geisonbiazus/blog/internal/adapters/tokenencoder"
	"github.com/geisonbiazus/blog/internal/adapters/transactionmanager"
	"github.com/geisonbiazus/blog/internal/adapters/userrepo"
	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/internal/core/shared"
	"github.com/geisonbiazus/blog/internal/ui/web"
	"github.com/geisonbiazus/blog/pkg/env"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Context struct {
	Env string

	Port         int
	TemplatePath string
	StaticPath   string
	PostPath     string
	BaseURL      string

	GitHubClientID     string
	GitHubClientSecret string

	AuthTokenSecret string

	PostgresURL string

	db        *sql.DB
	stateRepo auth.StateRepo
	userRepo  auth.UserRepo

	renderedPostCache shared.Cache[blog.RenderedPost]
}

func NewContext() *Context {
	return &Context{
		Env: env.GetString("ENV", "development"),

		Port:         env.GetInt("PORT", 3000),
		TemplatePath: env.GetString("TEMPLATE_PATH", filepath.Join("web", "template")),
		StaticPath:   env.GetString("STATIC_PATH", filepath.Join("web", "static")),
		PostPath:     env.GetString("POST_PATH", filepath.Join("posts")),
		BaseURL:      env.GetString("BASE_URL", "http://localhost:3000"),

		GitHubClientID:     env.GetString("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: env.GetString("GITHUB_CLIENT_SECRET", ""),

		AuthTokenSecret: env.GetString("AUTH_TOKEN_SECRET", ""),

		PostgresURL: env.GetString("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/blog?sslmode=disable"),
	}
}

// UI

func (c *Context) WebServer() *web.Server {
	return web.NewServer(c.Port, c.Router(), c.Logger())
}

func (c *Context) Router() http.Handler {
	return web.NewRouter(c.TemplatePath, c.StaticPath, c.UseCases(), c.BaseURL)
}

// Use cases

func (c *Context) UseCases() *web.UseCases {
	return &web.UseCases{
		ViewPost:      c.ViewPostUseCase(),
		ListPosts:     c.ListPostsUseCase(),
		RequestOAuth2: c.RequestOAuth2UseCase(),
		ConfirmOAuth2: c.ConfirmOAuth2UseCase(),
	}
}

func (c *Context) ViewPostUseCase() *blog.ViewPostUseCase {
	return blog.NewViewPostUseCase(c.PostRepo(), c.Renderer(), c.RenderedPostCache())
}

func (c *Context) ListPostsUseCase() *blog.ListPostsUseCase {
	return blog.NewListPostsUseCase(c.PostRepo(), c.Renderer())
}

func (c *Context) RequestOAuth2UseCase() *auth.RequestOAuth2UseCase {
	return auth.NewRequestOAuth2UseCase(c.OAuth2Provider(), c.IDGenerator(), c.StateRepo())
}

func (c *Context) ConfirmOAuth2UseCase() *auth.ConfirmOAuth2UseCase {
	return auth.NewConfirmOAuth2UseCase(c.OAuth2Provider(), c.StateRepo(), c.UserRepo(), c.IDGenerator(), c.TokenEncoder(), c.TransactionManager())
}

// Adapters

func (c *Context) DB() *sql.DB {
	if c.db == nil {
		db, err := sql.Open("pgx", c.PostgresURL)
		if err != nil {
			panic(err)
		}
		c.db = db
	}
	return c.db
}

func (c *Context) TransactionManager() shared.TransactionManager {
	if c.isTest() {
		return transactionmanager.NewFakeTransactionManager()
	}
	return transactionmanager.NewPostgresTransactionManager(c.DB())
}

func (c *Context) PostRepo() blog.PostRepo {
	return postrepo.NewFileSystemPostRepo(c.PostPath)
}

func (c *Context) Renderer() blog.Renderer {
	return renderer.NewGoldmarkRenderer()
}

func (c *Context) RenderedPostCache() shared.Cache[blog.RenderedPost] {
	if c.renderedPostCache == nil {
		c.renderedPostCache = cache.NewMemoryCache[blog.RenderedPost]()
	}
	return c.renderedPostCache
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

func (c *Context) IDGenerator() auth.IDGenerator {
	return idgenerator.NewUUIDGenerator()
}

func (c *Context) StateRepo() auth.StateRepo {
	if c.stateRepo == nil {
		c.stateRepo = staterepo.NewMemoryStateRepo()
	}
	return c.stateRepo
}

func (c *Context) UserRepo() auth.UserRepo {
	if c.userRepo == nil {
		if c.isTest() {
			c.userRepo = userrepo.NewMemoryUserRepo()
		} else {
			c.userRepo = userrepo.NewPostgresUserRepo(c.DB())
		}
	}
	return c.userRepo
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
