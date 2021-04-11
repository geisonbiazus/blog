package testhelper

import "github.com/geisonbiazus/blog/internal/core/posts"

type ViewPostUseCaseSpy struct {
	ReceivedPath string
	RenderedPost posts.RenderedPost
	Error        error
}

func (u *ViewPostUseCaseSpy) Run(path string) (posts.RenderedPost, error) {
	u.ReceivedPath = path
	return u.RenderedPost, u.Error
}
