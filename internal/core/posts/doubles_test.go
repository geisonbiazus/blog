package posts_test

import (
	"github.com/geisonbiazus/blog/internal/core/posts"
)

type PostRepoSpy struct {
	ReceivedPath string
	ReturnPost   posts.Post
	ReturnError  error
}

func (f *PostRepoSpy) GetPostByPath(path string) (posts.Post, error) {
	f.ReceivedPath = path
	return f.ReturnPost, f.ReturnError
}

type RendererSpy struct {
	ReceivedContent       string
	ReturnError           error
	ReturnRenderedContent string
}

func (r *RendererSpy) Render(content string) (string, error) {
	r.ReceivedContent = content
	return r.ReturnRenderedContent, r.ReturnError
}
