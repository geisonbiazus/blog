package app

import (
	"net/http"

	"github.com/geisonbiazus/blog/internal/adapters/postrepo/filesystem"
	"github.com/geisonbiazus/blog/internal/adapters/renderer/goldmark"
	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/internal/ui/web"
	"github.com/geisonbiazus/blog/pkg/env"
)

type Context struct {
	Port         int
	TemplatePath string
	PostPath     string
}

func NewContext() *Context {
	return &Context{
		Port:         env.GetInt("PORT", 3000),
		TemplatePath: env.GetString("TEMPLATE_PATH", "web/template"),
		PostPath:     env.GetString("POST_PATH", "posts"),
	}
}

func (c *Context) WebServer() *web.Server {
	return web.NewServer(c.Port, c.Router())
}

func (c *Context) Router() http.Handler {
	router, err := web.NewRouter(c.TemplatePath, c.UseCases())

	if err != nil {
		panic(err)
	}

	return router
}

func (c *Context) UseCases() *web.UseCases {
	return &web.UseCases{
		ViewPost: c.ViewPostUseCase(),
	}
}

func (c *Context) ViewPostUseCase() *posts.ViewPostUseCase {
	return posts.NewVewPostUseCase(c.PostRepo(), c.Renderer())
}

func (c *Context) PostRepo() *filesystem.PostRepo {
	return filesystem.NewPostRepo(c.PostPath)
}

func (c *Context) Renderer() *goldmark.GoldmarkRenderer {
	return goldmark.NewGoldmarkRenderer()
}
