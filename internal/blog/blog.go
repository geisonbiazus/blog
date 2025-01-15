package blog

import (
	"path/filepath"

	"github.com/geisonbiazus/blog/internal/blog/adapters/postrepo"
	"github.com/geisonbiazus/blog/internal/blog/adapters/renderer"
	"github.com/geisonbiazus/blog/internal/blog/entities"
	"github.com/geisonbiazus/blog/internal/blog/ports"
	"github.com/geisonbiazus/blog/internal/blog/usecases"
	"github.com/geisonbiazus/blog/pkg/caching"
	"github.com/geisonbiazus/blog/pkg/env"
)

type Post = entities.Post
type RenderedPost = entities.RenderedPost
type PostRepo = ports.PostRepo
type Renderer = ports.Renderer

var ErrPostNotFound = entities.ErrPostNotFound

type Context struct {
	postPath string
	cache    func() caching.Cache
}

func NewContext(cache func() caching.Cache) *Context {
	return &Context{
		postPath: env.GetString("POST_PATH", filepath.Join("posts")),
		cache:    cache,
	}
}

// Use cases

func (c *Context) ViewPostUseCase() *usecases.ViewPostUseCase {
	return usecases.NewViewPostUseCase(c.postRepo(), c.renderer(), c.cache())
}

func (c *Context) ListPostsUseCase() *usecases.ListPostsUseCase {
	return usecases.NewListPostsUseCase(c.postRepo(), c.renderer(), c.cache())
}

// Adapters

func (c *Context) postRepo() PostRepo {
	return postrepo.NewFileSystemPostRepo(c.postPath)
}

func (c *Context) renderer() Renderer {
	return renderer.NewGoldmarkRenderer()
}
