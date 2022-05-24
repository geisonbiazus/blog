package blog

import "github.com/geisonbiazus/blog/internal/core/shared"

type ViewPostUseCase struct {
	postRepo PostRepo
	renderer Renderer
	cache    shared.Cache[RenderedPost]
}

func NewViewPostUseCase(
	postRepo PostRepo,
	renderer Renderer,
	cache shared.Cache[RenderedPost],
) *ViewPostUseCase {
	return &ViewPostUseCase{
		postRepo: postRepo,
		renderer: renderer,
		cache:    cache,
	}
}

func (u *ViewPostUseCase) Run(path string) (RenderedPost, error) {
	return u.cache.Do(path, func() (RenderedPost, error) {
		return u.run(path)
	}, shared.NeverExpire)
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
