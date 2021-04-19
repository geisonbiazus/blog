package posts_test

import (
	"errors"
	"testing"

	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type listPostsUseCaseFixture struct {
	usecase *posts.ListPostsUseCase
	repo    *PostRepoSpy
}

func TestTestListPostsUseCase(t *testing.T) {
	setup := func() *listPostsUseCaseFixture {
		repo := NewPostRepoSpy()
		usecase := posts.NewListPostsUseCase(repo)
		return &listPostsUseCaseFixture{
			usecase: usecase,
			repo:    repo,
		}
	}

	t.Run("Given not posts, it returns an empty slice", func(t *testing.T) {
		f := setup()

		result, err := f.usecase.Run()

		assert.DeepEqual(t, []posts.Post{}, result)
		assert.Nil(t, err)
	})

	t.Run("Given some posts, it returns them", func(t *testing.T) {
		f := setup()

		postList := []posts.Post{newPost()}
		f.repo.ReturnPosts = postList

		result, err := f.usecase.Run()

		assert.DeepEqual(t, postList, result)
		assert.Nil(t, err)
	})

	t.Run("Given an error, it returns the error", func(t *testing.T) {
		f := setup()

		f.repo.ReturnError = errors.New("Error")

		result, err := f.usecase.Run()

		assert.DeepEqual(t, []posts.Post{}, result)
		assert.Equal(t, f.repo.ReturnError, err)
	})
}
