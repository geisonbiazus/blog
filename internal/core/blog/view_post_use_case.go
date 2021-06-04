package blog

type ViewPostUseCase struct {
	postRepo PostRepo
	renderer Renderer
}

func NewViewPostUseCase(postRepo PostRepo, renderer Renderer) *ViewPostUseCase {
	return &ViewPostUseCase{postRepo: postRepo, renderer: renderer}
}

func (u *ViewPostUseCase) Run(path string) (RenderedPost, error) {
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
