package usecases_test

import "github.com/geisonbiazus/blog/internal/blog/entities"

type PostRepoSpy struct {
	ReceivedPath string
	ReturnPost   entities.Post
	ReturnPosts  []entities.Post
	ReturnError  error
}

func NewPostRepoSpy() *PostRepoSpy {
	return &PostRepoSpy{ReturnPosts: []entities.Post{}}
}

func (r *PostRepoSpy) GetPostByPath(path string) (entities.Post, error) {
	r.ReceivedPath = path
	return r.ReturnPost, r.ReturnError
}

func (r *PostRepoSpy) GetAllPosts() ([]entities.Post, error) {
	return r.ReturnPosts, r.ReturnError
}

type RendererSpy struct {
	ReceivedContent       string
	ReturnError           error
	ReturnRenderedContent string
}

func NewRendererSpy() *RendererSpy {
	return &RendererSpy{}
}

func (r *RendererSpy) Render(content string) (string, error) {
	r.ReceivedContent = content
	return r.ReturnRenderedContent, r.ReturnError
}
