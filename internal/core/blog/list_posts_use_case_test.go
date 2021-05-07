package blog_test

import (
	"errors"
	"testing"

	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type listPostsUseCaseFixture struct {
	usecase *blog.ListPostsUseCase
	repo    *PostRepoSpy
}

func TestTestListPostsUseCase(t *testing.T) {
	setup := func() *listPostsUseCaseFixture {
		repo := NewPostRepoSpy()
		usecase := blog.NewListPostsUseCase(repo)
		return &listPostsUseCaseFixture{
			usecase: usecase,
			repo:    repo,
		}
	}

	t.Run("Given no post exists, it returns an empty slice", func(t *testing.T) {
		f := setup()

		posts, err := f.usecase.Run()

		assert.DeepEqual(t, []blog.Post{}, posts)
		assert.Nil(t, err)
	})

	t.Run("Given some posts, it returns them", func(t *testing.T) {
		f := setup()

		posts := []blog.Post{newPost()}
		f.repo.ReturnPosts = posts

		result, err := f.usecase.Run()

		assert.DeepEqual(t, posts, result)
		assert.Nil(t, err)
	})

	t.Run("Given an error, it returns the error", func(t *testing.T) {
		f := setup()

		f.repo.ReturnError = errors.New("Error")

		result, err := f.usecase.Run()

		assert.DeepEqual(t, []blog.Post{}, result)
		assert.Equal(t, f.repo.ReturnError, err)
	})
}
