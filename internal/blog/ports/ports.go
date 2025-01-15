package ports

import "github.com/geisonbiazus/blog/internal/blog/entities"

type PostRepo interface {
	GetPostByPath(path string) (entities.Post, error)
	GetAllPosts() ([]entities.Post, error)
}

type Renderer interface {
	Render(content string) (string, error)
}
