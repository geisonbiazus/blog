package ports

import "github.com/geisonbiazus/blog/internal/blog/entities"

// Driving ports

type ListPostsUseCase interface {
	Run() ([]entities.RenderedPost, error)
}

type ViewPostUseCase interface {
	Run(path string) (entities.RenderedPost, error)
}

// Driven ports

type PostRepo interface {
	GetPostByPath(path string) (entities.Post, error)
	GetAllPosts() ([]entities.Post, error)
}

type Renderer interface {
	Render(content string) (string, error)
}
