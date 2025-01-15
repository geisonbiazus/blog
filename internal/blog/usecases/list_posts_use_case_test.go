package usecases_test

import (
	"errors"
	"testing"

	"github.com/geisonbiazus/blog/internal/blog/entities"
	"github.com/geisonbiazus/blog/internal/blog/usecases"
	"github.com/geisonbiazus/blog/pkg/caching"
	"github.com/stretchr/testify/assert"
)

type listPostsUseCaseFixture struct {
	usecase  *usecases.ListPostsUseCase
	repo     *PostRepoSpy
	renderer *RendererSpy
}

func TestTestListPostsUseCase(t *testing.T) {
	setup := func() *listPostsUseCaseFixture {
		repo := NewPostRepoSpy()
		renderer := NewRendererSpy()
		cache := caching.NewMemoryCache()
		usecase := usecases.NewListPostsUseCase(repo, renderer, cache)
		return &listPostsUseCaseFixture{
			usecase:  usecase,
			repo:     repo,
			renderer: renderer,
		}
	}

	t.Run("Given no post exists, it returns an empty slice", func(t *testing.T) {
		f := setup()

		posts, err := f.usecase.Run()

		assert.Equal(t, []entities.RenderedPost{}, posts)
		assert.Nil(t, err)
	})

	t.Run("Given some posts, it renders and returns them", func(t *testing.T) {
		f := setup()

		post := newPost()
		renderedPost := entities.RenderedPost{Post: post, HTML: "Rendered post"}
		posts := []entities.Post{post}
		renderedPosts := []entities.RenderedPost{renderedPost}

		f.repo.ReturnPosts = posts
		f.renderer.ReturnRenderedContent = renderedPost.HTML

		result, err := f.usecase.Run()

		assert.Equal(t, renderedPosts, result)
		assert.Nil(t, err)
	})

	t.Run("Given an error is returned form the repo, it returns the error", func(t *testing.T) {
		f := setup()

		f.repo.ReturnError = errors.New("Repo error")

		result, err := f.usecase.Run()

		assert.Equal(t, []entities.RenderedPost{}, result)
		assert.Equal(t, f.repo.ReturnError, err)
	})

	t.Run("Given an error is returned form the renderer, it returns the error", func(t *testing.T) {
		f := setup()

		f.repo.ReturnPosts = []entities.Post{newPost()}
		f.renderer.ReturnError = errors.New("Renderer error")

		result, err := f.usecase.Run()

		assert.Equal(t, []entities.RenderedPost{}, result)
		assert.Equal(t, f.renderer.ReturnError, err)
	})
}
