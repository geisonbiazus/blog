package filesystem_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/postrepo/filesystem"
	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/stretchr/testify/assert"
)

const postPath = "fixtures/posts"
const emptyFolder = "fixtures/empty_folder"
const invalidPath = "fixtures/invalid-path"
const pathWithInvalidPost = "fixtures/posts_with_invalid"
const pathWithDifferentFiles = "fixtures/posts_with_other_files"

func TestPostRepo(t *testing.T) {
	t.Run("GetPostByPath", func(t *testing.T) {
		t.Run("It returns error when post is not found", func(t *testing.T) {
			repo := filesystem.NewPostRepo(postPath)

			post, err := repo.GetPostByPath("wrong_path")

			assert.Equal(t, blog.Post{}, post)
			assert.Equal(t, blog.ErrPostNotFound, err)
		})

		t.Run("It returns parsed post when the file exists", func(t *testing.T) {
			repo := filesystem.NewPostRepo(postPath)

			post, err := repo.GetPostByPath("test-post-1")

			assert.Nil(t, err)
			assert.Equal(t, testPost1, post)
		})
	})

	t.Run("GetAllPosts()", func(t *testing.T) {
		t.Run("Given a path with no post files, it returns empty", func(t *testing.T) {
			repo := filesystem.NewPostRepo(emptyFolder)
			posts, err := repo.GetAllPosts()

			assert.Equal(t, []blog.Post{}, posts)
			assert.Nil(t, err)
		})

		t.Run("Given a non-existent path, it returns error", func(t *testing.T) {
			repo := filesystem.NewPostRepo(invalidPath)
			posts, err := repo.GetAllPosts()

			assert.Equal(t, []blog.Post{}, posts)
			assert.NotNil(t, err)
		})

		t.Run("Given a path with post files, it returns all posts sorted by descending date", func(t *testing.T) {
			repo := filesystem.NewPostRepo(postPath)
			expectedPosts := []blog.Post{testPost1, testPost3, testPost2}

			actualPosts, err := repo.GetAllPosts()

			assert.Equal(t, expectedPosts, actualPosts)
			assert.Nil(t, err)
		})

		t.Run("Given an invalid post in the folder, it ignores the invalid and returns the rest", func(t *testing.T) {
			repo := filesystem.NewPostRepo(pathWithInvalidPost)
			expectedPosts := []blog.Post{testPost1, testPost2}

			actualPosts, err := repo.GetAllPosts()

			assert.Equal(t, expectedPosts, actualPosts)
			assert.Nil(t, err)
		})

		t.Run("Given a path other types of files, it ignores the other files", func(t *testing.T) {
			repo := filesystem.NewPostRepo(pathWithDifferentFiles)
			expectedPosts := []blog.Post{testPost1}

			actualPosts, err := repo.GetAllPosts()

			assert.Equal(t, expectedPosts, actualPosts)
			assert.Nil(t, err)
		})
	})
}

var testPost1 = blog.Post{
	Title:       "Test Post 1",
	Author:      "Geison Biazus",
	Path:        "test-post-1",
	Description: "Description of post 1",
	ImagePath:   "/post-image-1.png",
	Time:        toTime("2021-04-05T18:47:00Z"),
	Markdown: "" +
		"## Subtitle\n" +
		"\n" +
		"Content\n",
}

var testPost2 = blog.Post{
	Title:       "Test Post 2",
	Author:      "Geison Biazus",
	Path:        "test-post-2",
	Description: "Description of post 2",
	ImagePath:   "/post-image-2.png",
	Time:        toTime("2021-04-04T14:33:00Z"),
	Markdown:    "",
}

var testPost3 = blog.Post{
	Title:       "Test Post 3",
	Author:      "Geison Biazus",
	Path:        "test-post-3",
	Description: "Description of post 3",
	ImagePath:   "/post-image-3.png",
	Time:        toTime("2021-04-05T18:40:00Z"),
	Markdown:    "",
}
