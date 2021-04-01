package posts_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type viewPostUseCaseFixture struct {
	usecase  *posts.ViewPostUseCase
	post     posts.Post
	repo     *FakePostRepo
	renderer *FakeRenderer
}

func TestViewPostUseCase(t *testing.T) {
	setup := func() *viewPostUseCaseFixture {
		repo := &FakePostRepo{}
		renderer := &FakeRenderer{}
		usecase := posts.NewVewPostUseCase(repo, renderer)

		postTime, _ := time.Parse(time.RFC3339, "2021-04-03T00:00:00+00:00")
		post := posts.Post{
			Title:   "Title",
			Authors: []string{"Author"},
			Time:    postTime,
			Path:    "path",
			Content: "content",
		}
		repo.AddPost(post)

		return &viewPostUseCaseFixture{
			usecase:  usecase,
			post:     post,
			repo:     repo,
			renderer: renderer,
		}
	}

	t.Run("It returns error when post is not found", func(t *testing.T) {
		f := setup()

		body, err := f.usecase.Run("wrong-path")

		assert.DeepEqual(t, posts.RenderedPost{}, body)
		assert.Equal(t, posts.ErrPostNotFound, err)
	})

	t.Run("It returns the rendered post when post is found", func(t *testing.T) {
		f := setup()

		rennderedPost, err := f.usecase.Run(f.post.Path)

		renderedContent, _ := f.renderer.Render(f.post.Content)
		assert.DeepEqual(t, rennderedPost, posts.RenderedPost{
			Title:   f.post.Title,
			Authors: f.post.Authors,
			Time:    f.post.Time,
			Content: renderedContent,
		})
		assert.Nil(t, err)
	})

	t.Run("It returns error when fails to render", func(t *testing.T) {
		f := setup()

		f.renderer.RenderError = errors.New("render error")

		rennderedPost, err := f.usecase.Run(f.post.Path)

		assert.DeepEqual(t, rennderedPost, posts.RenderedPost{})
		assert.Equal(t, f.renderer.RenderError, err)
	})
}

type FakePostRepo struct {
	posts []posts.Post
}

func (f *FakePostRepo) AddPost(post posts.Post) {
	f.posts = append(f.posts, post)
}

func (f *FakePostRepo) GetPostByPath(path string) (posts.Post, error) {
	for _, post := range f.posts {
		if post.Path == path {
			return post, nil
		}
	}
	return posts.Post{}, posts.ErrPostNotFound
}

type FakeRenderer struct {
	RenderError error
}

func (r *FakeRenderer) Render(content string) (string, error) {
	if r.RenderError != nil {
		return "", r.RenderError
	}

	return fmt.Sprintf("Rendered %s", content), nil
}
