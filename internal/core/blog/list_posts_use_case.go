package blog

type ListPostsUseCase struct {
	postRepo PostRepo
	renderer Renderer
}

func NewListPostsUseCase(postRepo PostRepo, renderer Renderer) *ListPostsUseCase {
	return &ListPostsUseCase{postRepo: postRepo, renderer: renderer}
}

func (u *ListPostsUseCase) Run() ([]RenderedPost, error) {
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
