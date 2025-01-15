package usecases

import (
	"github.com/geisonbiazus/blog/internal/blog/entities"
	"github.com/geisonbiazus/blog/internal/blog/ports"
	"github.com/geisonbiazus/blog/pkg/caching"
)

type ViewPostUseCase struct {
	postRepo ports.PostRepo
	renderer ports.Renderer
	cache    caching.Cache
}

func NewViewPostUseCase(
	postRepo ports.PostRepo,
	renderer ports.Renderer,
	cache caching.Cache,
) *ViewPostUseCase {
	return &ViewPostUseCase{
		postRepo: postRepo,
		renderer: renderer,
		cache:    cache,
	}
}

func (u *ViewPostUseCase) Run(path string) (entities.RenderedPost, error) {
	result, err := u.cache.Do(path, func() (interface{}, error) {
		return u.run(path)
	}, caching.NeverExpire)

	return result.(entities.RenderedPost), err
}

func (u *ViewPostUseCase) run(path string) (entities.RenderedPost, error) {
	post, err := u.postRepo.GetPostByPath(path)

	if err != nil {
		return entities.RenderedPost{}, err
	}

	return u.renderPost(post)
}

func (u *ViewPostUseCase) renderPost(post entities.Post) (entities.RenderedPost, error) {
	renderedContent, err := u.renderer.Render(post.Markdown)

	if err != nil {
		return entities.RenderedPost{}, err
	}

	return entities.RenderedPost{
		Post: post,
		HTML: renderedContent,
	}, nil
}
