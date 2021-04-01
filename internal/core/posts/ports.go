package posts

import "errors"

type PostRepo interface {
	GetPostByPath(path string) (Post, error)
}

var ErrPostNotFound = errors.New("post not found")

type Renderer interface {
	Render(content string) (string, error)
}
