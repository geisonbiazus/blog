package usecases

import (
	"github.com/geisonbiazus/blog/internal/blog/entities"
	"github.com/geisonbiazus/blog/internal/blog/ports"
	"github.com/geisonbiazus/blog/pkg/caching"
)

type ListPostsUseCase struct {
	postRepo ports.PostRepo
	renderer ports.Renderer
	cache    caching.Cache
}

func NewListPostsUseCase(
	postRepo ports.PostRepo,
	renderer ports.Renderer,
	cache caching.Cache,
) *ListPostsUseCase {
	return &ListPostsUseCase{
		postRepo: postRepo,
		renderer: renderer,
		cache:    cache,
	}
}

const cacheKey = "all-posts"

func (u *ListPostsUseCase) Run() ([]entities.RenderedPost, error) {
	result, err := u.cache.Do(cacheKey, func() (interface{}, error) {
		return u.run()
	}, caching.NeverExpire)

	return result.([]entities.RenderedPost), err
}

func (u *ListPostsUseCase) run() ([]entities.RenderedPost, error) {
	posts, err := u.postRepo.GetAllPosts()

	if err != nil {
		return []entities.RenderedPost{}, err
	}

	return u.renderPosts(posts)
}

func (u *ListPostsUseCase) renderPosts(posts []entities.Post) ([]entities.RenderedPost, error) {
	renderedPosts := []entities.RenderedPost{}

	for _, post := range posts {
		html, err := u.renderer.Render(post.Markdown)

		if err != nil {
			return []entities.RenderedPost{}, err
		}

		renderedPosts = append(renderedPosts, entities.RenderedPost{Post: post, HTML: html})
	}

	return renderedPosts, nil
}
