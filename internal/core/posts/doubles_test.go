package posts_test

import (
	"github.com/geisonbiazus/blog/internal/core/posts"
)

type PostRepoSpy struct {
	ReceivedPath string
	ReturnPost   posts.Post
	ReturnPosts  []posts.Post
	ReturnError  error
}

func NewPostRepoSpy() *PostRepoSpy {
	return &PostRepoSpy{ReturnPosts: []posts.Post{}}
}

func (r *PostRepoSpy) GetPostByPath(path string) (posts.Post, error) {
	r.ReceivedPath = path
	return r.ReturnPost, r.ReturnError
}

func (r *PostRepoSpy) GetAllPosts() ([]posts.Post, error) {
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
