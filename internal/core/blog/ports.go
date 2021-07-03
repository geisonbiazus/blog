package blog

type PostRepo interface {
	GetPostByPath(path string) (Post, error)
	GetAllPosts() ([]Post, error)
}

type Renderer interface {
	Render(content string) (string, error)
}
