package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/geisonbiazus/blog/internal/adapters/idgenerator/uuid"
	"github.com/geisonbiazus/blog/internal/adapters/oauth2provider/fake"
	"github.com/geisonbiazus/blog/internal/adapters/oauth2provider/github"
	"github.com/geisonbiazus/blog/internal/adapters/postrepo/filesystem"
	"github.com/geisonbiazus/blog/internal/adapters/renderer/goldmark"
	staterepo "github.com/geisonbiazus/blog/internal/adapters/staterepo/memory"
	"github.com/geisonbiazus/blog/internal/adapters/tokenencoder/jwt"
	userrepo_memory "github.com/geisonbiazus/blog/internal/adapters/userrepo/memory"
	userrepo_postgres "github.com/geisonbiazus/blog/internal/adapters/userrepo/postgres"
	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/internal/ui/web"
	"github.com/geisonbiazus/blog/pkg/env"
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
	return blog.NewViewPostUseCase(c.PostRepo(), c.Renderer())
}

func (c *Context) ListPostsUseCase() *blog.ListPostsUseCase {
	return blog.NewListPostsUseCase(c.PostRepo(), c.Renderer())
}

func (c *Context) RequestOAuth2UseCase() *auth.RequestOAuth2UseCase {
	return auth.NewRequestOAuth2UseCase(c.OAuth2Provider(), c.IDGenerator(), c.StateRepo())
}

func (c *Context) ConfirmOAuth2UseCase() *auth.ConfirmOAuth2UseCase {
	return auth.NewConfirmOAuth2UseCase(c.OAuth2Provider(), c.StateRepo(), c.UserRepo(), c.IDGenerator(), c.TokenManager())
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

func (c *Context) PostRepo() *filesystem.PostRepo {
	return filesystem.NewPostRepo(c.PostPath)
}

func (c *Context) Renderer() *goldmark.Renderer {
	return goldmark.NewRenderer()
}

func (c *Context) OAuth2Provider() auth.OAuth2Provider {
	if c.Env == "test" {
		return c.FakeOAuth2Provider()
	}
	return c.GithubOAuth2Provider()
}

func (c *Context) GithubOAuth2Provider() *github.Provider {
	return github.NewProvider(c.GitHubClientID, c.GitHubClientSecret)
}

func (c *Context) FakeOAuth2Provider() *fake.Provider {
	return fake.NewProvider(c.BaseURL)
}

func (c *Context) IDGenerator() *uuid.Generator {
	return uuid.NewGenerator()
}

func (c *Context) StateRepo() auth.StateRepo {
	if c.stateRepo == nil {
		c.stateRepo = staterepo.NewStateRepo()
	}
	return c.stateRepo
}

func (c *Context) UserRepo() auth.UserRepo {
	if c.userRepo == nil {
		if c.Env == "test" {
			c.userRepo = userrepo_memory.NewUserRepo()
		} else {
			c.userRepo = userrepo_postgres.NewUserRepo(c.DB())
		}
	}
	return c.userRepo
}

func (c *Context) TokenManager() *jwt.TokenEncoder {
	return jwt.NewTokenEncoder(c.AuthTokenSecret)
}

func (c *Context) Logger() *log.Logger {
	return log.New(os.Stdout, "web: ", log.Ldate|log.Ltime|log.LUTC)
}
