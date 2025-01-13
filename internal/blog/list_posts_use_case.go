package blog

import "github.com/geisonbiazus/blog/pkg/caching"

type ListPostsUseCase struct {
	postRepo PostRepo
	renderer Renderer
	cache    caching.Cache
}

func NewListPostsUseCase(
	postRepo PostRepo,
	renderer Renderer,
	cache caching.Cache,
) *ListPostsUseCase {
	return &ListPostsUseCase{
		postRepo: postRepo,
		renderer: renderer,
		cache:    cache,
	}
}

const cacheKey = "all-posts"

func (u *ListPostsUseCase) Run() ([]RenderedPost, error) {
	result, err := u.cache.Do(cacheKey, func() (interface{}, error) {
		return u.run()
	}, caching.NeverExpire)

	return result.([]RenderedPost), err
}

func (u *ListPostsUseCase) run() ([]RenderedPost, error) {
	posts, err := u.postRepo.GetAllPosts()

	if err != nil {
		return []RenderedPost{}, err
	}

	return u.renderPosts(posts)
}

func (u *ListPostsUseCase) renderPosts(posts []Post) ([]RenderedPost, error) {
	renderedPosts := []RenderedPost{}

	for _, post := range posts {
		html, err := u.renderer.Render(post.Markdown)

		if err != nil {
			return []RenderedPost{}, err
		}

		renderedPosts = append(renderedPosts, RenderedPost{Post: post, HTML: html})
	}

	return renderedPosts, nil
}
