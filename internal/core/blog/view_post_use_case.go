package blog

import "github.com/geisonbiazus/blog/internal/core/shared"

type ViewPostUseCase struct {
	postRepo PostRepo
	renderer Renderer
	cache    shared.Cache
}

func NewViewPostUseCase(
	postRepo PostRepo,
	renderer Renderer,
	cache shared.Cache,
) *ViewPostUseCase {
	return &ViewPostUseCase{
		postRepo: postRepo,
		renderer: renderer,
		cache:    cache,
	}
}

func (u *ViewPostUseCase) Run(path string) (RenderedPost, error) {
	result, err := u.cache.Do(path, func() (interface{}, error) {
		return u.run(path)
	}, shared.NeverExpire)

	return result.(RenderedPost), err
}

func (u *ViewPostUseCase) run(path string) (RenderedPost, error) {
	post, err := u.postRepo.GetPostByPath(path)

	if err != nil {
		return RenderedPost{}, err
	}

	return u.renderPost(post)
}

func (u *ViewPostUseCase) renderPost(post Post) (RenderedPost, error) {
	renderedContent, err := u.renderer.Render(post.Markdown)

	if err != nil {
		return RenderedPost{}, err
	}

	return RenderedPost{
		Post: post,
		HTML: renderedContent,
	}, nil
}
