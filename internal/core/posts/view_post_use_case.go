package posts

type ViewPostUseCase struct {
	postRepo PostRepo
	renderer Renderer
}

func NewVewPostUseCase(postRepo PostRepo, renderer Renderer) *ViewPostUseCase {
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
	renderedContent, err := u.renderer.Render(post.Content)

	if err != nil {
		return RenderedPost{}, err
	}

	return RenderedPost{
		Title:   post.Title,
		Authors: post.Authors,
		Time:    post.Time,
		Content: renderedContent,
	}, nil
}
