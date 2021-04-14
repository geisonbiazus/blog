package app

import (
	"net/http"

	"github.com/geisonbiazus/blog/internal/adapters/postrepo/filesystem"
	"github.com/geisonbiazus/blog/internal/adapters/renderer/goldmark"
	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/internal/ui/web"
)

type Context struct {
	Port         int
	TemplatePath string
	PostBasePath string
}

func NewContext() *Context {
	return &Context{
		Port:         3000,
		TemplatePath: "web/template",
		PostBasePath: "posts",
	}
}

func (c *Context) WebServer() *web.Server {
	return &web.Server{
		Port:   c.Port,
		Router: c.Router(),
	}
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
	return filesystem.NewPostRepo(c.PostBasePath)
}

func (c *Context) Renderer() *goldmark.GoldmarkRenderer {
	return goldmark.NewGoldmarkRenderer()
}
