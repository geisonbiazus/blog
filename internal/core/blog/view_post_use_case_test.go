package blog_test

import (
	"errors"
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/adapters/cache/memory"
	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type viewPostUseCaseFixture struct {
	usecase  *blog.ViewPostUseCase
	repo     *PostRepoSpy
	renderer *RendererSpy
}

func TestViewPostUseCase(t *testing.T) {
	setup := func() *viewPostUseCaseFixture {
		repo := NewPostRepoSpy()
		renderer := NewRendererSpy()
		cache := memory.NewCache()
		usecase := blog.NewViewPostUseCase(repo, renderer, cache)

		return &viewPostUseCaseFixture{
			usecase:  usecase,
			repo:     repo,
			renderer: renderer,
		}
	}

	t.Run("It returns error when post is not found", func(t *testing.T) {
		f := setup()

		f.repo.ReturnError = blog.ErrPostNotFound

		renderedPost, err := f.usecase.Run("path")

		assert.Equal(t, "path", f.repo.ReceivedPath)
		assert.DeepEqual(t, blog.RenderedPost{}, renderedPost)
		assert.Equal(t, blog.ErrPostNotFound, err)
	})

	t.Run("It returns the rendered post when post is found", func(t *testing.T) {
		f := setup()

		post := newPost()
		f.repo.ReturnPost = post
		f.renderer.ReturnRenderedContent = "Rendered content"

		rennderedPost, err := f.usecase.Run("path")

		assert.Equal(t, "path", f.repo.ReceivedPath)
		assert.Equal(t, post.Markdown, f.renderer.ReceivedContent)
		assert.Equal(t, rennderedPost, blog.RenderedPost{
			Post: post,
			HTML: "Rendered content",
		})
		assert.Nil(t, err)
	})

	t.Run("It returns error when post fails to render", func(t *testing.T) {
		f := setup()

		post := newPost()
		f.repo.ReturnPost = post
		f.renderer.ReturnError = errors.New("render error")

		rennderedPost, err := f.usecase.Run(post.Path)
		assert.Equal(t, rennderedPost, blog.RenderedPost{})
		assert.Equal(t, f.renderer.ReturnError, err)
	})
}

func newPost() blog.Post {
	postTime, _ := time.Parse(time.RFC3339, "2021-04-03T00:00:00+00:00")

	return blog.Post{
		Title:       "Title",
		Author:      "Author",
		Time:        postTime,
		Path:        "path",
		Description: "Description",
		ImagePath:   "/image.png",
		Markdown:    "content",
	}
}
