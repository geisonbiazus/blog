package post_repo_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/post_repo"
	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type fileSystemPostRepoFixture struct {
	repo *post_repo.FileSystemPostRepo
}

func TestFileSystemPostRepo(t *testing.T) {
	setup := func() *fileSystemPostRepoFixture {
		repo := post_repo.NewFileSystemPostRepo()

		return &fileSystemPostRepoFixture{
			repo: repo,
		}
	}

	t.Run("It returns error when post is not found", func(t *testing.T) {
		f := setup()

		post, err := f.repo.GetPostByPath("wrong_path")

		assert.DeepEqual(t, posts.Post{}, post)
		assert.Equal(t, posts.ErrPostNotFound, err)
	})
}
