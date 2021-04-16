package posts_test

import (
	"errors"
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type viewPostUseCaseFixture struct {
	usecase  *posts.ViewPostUseCase
	repo     *PostRepoSpy
	renderer *RendererSpy
}

func TestViewPostUseCase(t *testing.T) {
	setup := func() *viewPostUseCaseFixture {
		repo := NewPostRepoSpy()
		renderer := NewRendererSpy()
		usecase := posts.NewVewPostUseCase(repo, renderer)

		return &viewPostUseCaseFixture{
			usecase:  usecase,
			repo:     repo,
			renderer: renderer,
		}
	}

	t.Run("It returns error when post is not found", func(t *testing.T) {
		f := setup()

		f.repo.ReturnError = posts.ErrPostNotFound

		body, err := f.usecase.Run("path")

		assert.Equal(t, "path", f.repo.ReceivedPath)
		assert.DeepEqual(t, posts.RenderedPost{}, body)
		assert.Equal(t, posts.ErrPostNotFound, err)
	})

	t.Run("It returns the rendered post when post is found", func(t *testing.T) {
		f := setup()

		post := newPost()
		f.repo.ReturnPost = post
		f.renderer.ReturnRenderedContent = "Rendered content"

		rennderedPost, err := f.usecase.Run("path")

		assert.Equal(t, "path", f.repo.ReceivedPath)
		assert.Equal(t, post.Content, f.renderer.ReceivedContent)
		assert.DeepEqual(t, rennderedPost, posts.RenderedPost{
			Title:   post.Title,
			Author:  post.Author,
			Time:    post.Time,
			Content: "Rendered content",
		})
		assert.Nil(t, err)
	})

	t.Run("It returns error when fails to render", func(t *testing.T) {
		f := setup()

		f.renderer.ReturnError = errors.New("render error")

		rennderedPost, err := f.usecase.Run("path")
		assert.Equal(t, "path", f.repo.ReceivedPath)
		assert.DeepEqual(t, rennderedPost, posts.RenderedPost{})
		assert.Equal(t, f.renderer.ReturnError, err)
	})
}

func newPost() posts.Post {
	postTime, _ := time.Parse(time.RFC3339, "2021-04-03T00:00:00+00:00")

	return posts.Post{
		Title:   "Title",
		Author:  "Author",
		Time:    postTime,
		Path:    "path",
		Content: "content",
	}
}
