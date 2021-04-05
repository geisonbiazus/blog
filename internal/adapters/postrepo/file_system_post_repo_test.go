package postrepo_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/postrepo"
	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type fileSystemPostRepoFixture struct {
	repo *postrepo.FileSystemPostRepo
}

func TestFileSystemPostRepo(t *testing.T) {
	setup := func() *fileSystemPostRepoFixture {
		repo := postrepo.NewFileSystemPostRepo("./test_posts")

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

	t.Run("It returns parsed post when the file exists", func(t *testing.T) {
		f := setup()

		post, err := f.repo.GetPostByPath("test-post")

		assert.Nil(t, err)
		assert.DeepEqual(t, posts.Post{
			Title:  "Test Post",
			Author: "Geison Biazus",
			Path:   "test-post",
			Time:   toTime("2021-04-05T18:47:00Z"),
			Content: "" +
				"## Subtitle\n" +
				"\n" +
				"Content\n",
		}, post)
	})
}
