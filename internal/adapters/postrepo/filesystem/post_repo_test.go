package filesystem_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/postrepo/filesystem"
	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type postRepoFixture struct {
	repo *filesystem.PostRepo
}

func TestPostRepo(t *testing.T) {
	setup := func() *postRepoFixture {
		repo := filesystem.NewPostRepo("../../../../test/posts")

		return &postRepoFixture{
			repo: repo,
		}
	}

	t.Run("GetPostByPath", func(t *testing.T) {
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
			assert.DeepEqual(t, testPost1, post)
		})
	})

	t.Run("GetAllPosts()", func(t *testing.T) {
		t.Run("Given a path with no post files, it returns empty", func(t *testing.T) {
			repo := filesystem.NewPostRepo(".")
			postList, err := repo.GetAllPosts()

			assert.DeepEqual(t, []posts.Post{}, postList)
			assert.Nil(t, err)
		})

		t.Run("Given a non-existent path, it returns error", func(t *testing.T) {
			repo := filesystem.NewPostRepo("./invalid-path")
			postList, err := repo.GetAllPosts()

			assert.DeepEqual(t, []posts.Post{}, postList)
			assert.NotNil(t, err)
		})

		t.Run("Given a path with post files, it returns all posts sorted by descending date", func(t *testing.T) {
			f := setup()
			expectedPosts := []posts.Post{testPost1, testPost2, testPost3}

			actualPosts, err := f.repo.GetAllPosts()

			assert.DeepEqual(t, expectedPosts, actualPosts)
			assert.Nil(t, err)
		})
	})
}

var testPost1 = posts.Post{
	Title:  "Test Post",
	Author: "Geison Biazus",
	Path:   "test-post",
	Time:   toTime("2021-04-05T18:47:00Z"),
	Content: "" +
		"## Subtitle\n" +
		"\n" +
		"Content\n",
}

var testPost2 = posts.Post{
	Title:   "Test Post 3",
	Author:  "Geison Biazus",
	Path:    "test-post-3",
	Time:    toTime("2021-04-05T18:40:00Z"),
	Content: "",
}

var testPost3 = posts.Post{
	Title:   "Test Post 2",
	Author:  "Geison Biazus",
	Path:    "test-post-2",
	Time:    toTime("2021-04-04T14:33:00Z"),
	Content: "",
}
