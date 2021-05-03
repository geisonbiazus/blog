package app

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/geisonbiazus/blog/internal/adapters/postrepo/filesystem"
	"github.com/geisonbiazus/blog/internal/adapters/renderer/goldmark"
	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/internal/ui/web"
	"github.com/geisonbiazus/blog/pkg/env"
)

type Context struct {
	Port         int
	TemplatePath string
	StaticPath   string
	PostPath     string
}

func NewContext() *Context {
	return &Context{
		Port:         env.GetInt("PORT", 3000),
		TemplatePath: env.GetString("TEMPLATE_PATH", filepath.Join("web", "template")),
		StaticPath:   env.GetString("STATIC_PATH", filepath.Join("web", "static")),
		PostPath:     env.GetString("POST_PATH", filepath.Join("posts")),
	}
}

func (c *Context) WebServer() *web.Server {
	return web.NewServer(c.Port, c.Router(), c.Logger())
}

func (c *Context) Router() http.Handler {
	router, err := web.NewRouter(c.TemplatePath, c.StaticPath, c.UseCases())

	if err != nil {
		panic(err)
	}

	return router
}

func (c *Context) UseCases() *web.UseCases {
	return &web.UseCases{
		ViewPost:  c.ViewPostUseCase(),
		ListPosts: c.ListPostsUseCase(),
	}
}

func (c *Context) ViewPostUseCase() *blog.ViewPostUseCase {
	return blog.NewVewPostUseCase(c.PostRepo(), c.Renderer())
}

func (c *Context) ListPostsUseCase() *blog.ListPostsUseCase {
	return blog.NewListPostsUseCase(c.PostRepo())
}

func (c *Context) PostRepo() *filesystem.PostRepo {
	return filesystem.NewPostRepo(c.PostPath)
}

func (c *Context) Renderer() *goldmark.GoldmarkRenderer {
	return goldmark.NewGoldmarkRenderer()
}

func (c *Context) Logger() *log.Logger {
	return log.New(os.Stdout, "web: ", log.Ldate|log.Ltime|log.LUTC)
}
