package app

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/geisonbiazus/blog/internal/adapters/idgenerator/uuid"
	"github.com/geisonbiazus/blog/internal/adapters/oauth2provider/github"
	"github.com/geisonbiazus/blog/internal/adapters/postrepo/filesystem"
	"github.com/geisonbiazus/blog/internal/adapters/renderer/goldmark"
	"github.com/geisonbiazus/blog/internal/adapters/staterepo/memory"
	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/internal/ui/web"
	"github.com/geisonbiazus/blog/pkg/env"
)

type Context struct {
	Port         int
	TemplatePath string
	StaticPath   string
	PostPath     string
	BaseURL      string

	GitHubClientID     string
	GitHubClientSecret string
}

func NewContext() *Context {
	return &Context{
		Port:         env.GetInt("PORT", 3000),
		TemplatePath: env.GetString("TEMPLATE_PATH", filepath.Join("web", "template")),
		StaticPath:   env.GetString("STATIC_PATH", filepath.Join("web", "static")),
		PostPath:     env.GetString("POST_PATH", filepath.Join("posts")),
		BaseURL:      env.GetString("BASE_URL", "http://localhost:3000"),

		GitHubClientID:     env.GetString("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: env.GetString("GITHUB_CLIENT_SECRET", ""),
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
		RequestOauth2: c.RequestOauth2UseCase(),
	}
}

func (c *Context) ViewPostUseCase() *blog.ViewPostUseCase {
	return blog.NewViewPostUseCase(c.PostRepo(), c.Renderer())
}

func (c *Context) ListPostsUseCase() *blog.ListPostsUseCase {
	return blog.NewListPostsUseCase(c.PostRepo(), c.Renderer())
}

func (c *Context) RequestOauth2UseCase() *auth.RequestOauth2UseCase {
	return auth.NewRequestOauth2UseCase(c.Oauth2Provider(), c.IDGenerator(), c.StateRepo())
}

// Adapters

func (c *Context) PostRepo() *filesystem.PostRepo {
	return filesystem.NewPostRepo(c.PostPath)
}

func (c *Context) Renderer() *goldmark.Renderer {
	return goldmark.NewRenderer()
}

func (c *Context) Logger() *log.Logger {
	return log.New(os.Stdout, "web: ", log.Ldate|log.Ltime|log.LUTC)
}

func (c *Context) Oauth2Provider() *github.Provider {
	return github.NewProvider(c.GitHubClientID, c.GitHubClientSecret)
}

func (c *Context) IDGenerator() *uuid.Generator {
	return uuid.NewGenerator()
}

func (c *Context) StateRepo() *memory.InMemoryStateRepo {
	return memory.NewInMemoryStateRepo()
}
