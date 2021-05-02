package blog_test

import "github.com/geisonbiazus/blog/internal/core/blog"

type PostRepoSpy struct {
	ReceivedPath string
	ReturnPost   blog.Post
	ReturnPosts  []blog.Post
	ReturnError  error
}

func NewPostRepoSpy() *PostRepoSpy {
	return &PostRepoSpy{ReturnPosts: []blog.Post{}}
}

func (r *PostRepoSpy) GetPostByPath(path string) (blog.Post, error) {
	r.ReceivedPath = path
	return r.ReturnPost, r.ReturnError
}

func (r *PostRepoSpy) GetAllPosts() ([]blog.Post, error) {
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
