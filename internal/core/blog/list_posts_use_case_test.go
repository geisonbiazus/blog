package blog_test

import (
	"errors"
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/cache"
	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type listPostsUseCaseFixture struct {
	usecase  *blog.ListPostsUseCase
	repo     *PostRepoSpy
	renderer *RendererSpy
}

func TestTestListPostsUseCase(t *testing.T) {
	setup := func() *listPostsUseCaseFixture {
		repo := NewPostRepoSpy()
		renderer := NewRendererSpy()
		cache := cache.NewMemoryCache()
		usecase := blog.NewListPostsUseCase(repo, renderer, cache)
		return &listPostsUseCaseFixture{
			usecase:  usecase,
			repo:     repo,
			renderer: renderer,
		}
	}

	t.Run("Given no post exists, it returns an empty slice", func(t *testing.T) {
		f := setup()

		posts, err := f.usecase.Run()

		assert.DeepEqual(t, []blog.RenderedPost{}, posts)
		assert.Nil(t, err)
	})

	t.Run("Given some posts, it renders and returns them", func(t *testing.T) {
		f := setup()

		post := newPost()
		renderedPost := blog.RenderedPost{Post: post, HTML: "Rendered post"}
		posts := []blog.Post{post}
		renderedPosts := []blog.RenderedPost{renderedPost}

		f.repo.ReturnPosts = posts
		f.renderer.ReturnRenderedContent = renderedPost.HTML

		result, err := f.usecase.Run()

		assert.DeepEqual(t, renderedPosts, result)
		assert.Nil(t, err)
	})

	t.Run("Given an error is returned form the repo, it returns the error", func(t *testing.T) {
		f := setup()

		f.repo.ReturnError = errors.New("Repo error")

		result, err := f.usecase.Run()

		assert.DeepEqual(t, []blog.RenderedPost{}, result)
		assert.Equal(t, f.repo.ReturnError, err)
	})

	t.Run("Given an error is returned form the renderer, it returns the error", func(t *testing.T) {
		f := setup()

		f.repo.ReturnPosts = []blog.Post{newPost()}
		f.renderer.ReturnError = errors.New("Renderer error")

		result, err := f.usecase.Run()

		assert.DeepEqual(t, []blog.RenderedPost{}, result)
		assert.Equal(t, f.renderer.ReturnError, err)
	})
}
