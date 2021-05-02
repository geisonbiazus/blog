package filesystem_test

import (
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/adapters/postrepo/filesystem"
	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestParseFileContent(t *testing.T) {
	t.Run("It parses content header into a Post", func(t *testing.T) {
		assertParsedContent(t, "title: Post Title\n--\n", blog.Post{Title: "Post Title"})
		assertParsedContent(t, "author: Author Name\n--\n", blog.Post{Author: "Author Name"})
		assertParsedContent(t, "time: 2021-04-04 22:00\n--\n", blog.Post{Time: toTime("2021-04-04T22:00:00Z")})
		assertParsedContent(t, ""+
			"title: Post Title\n"+
			"author: Author Name\n"+
			"time: 2021-04-04 22:00\n"+
			"--\n",
			blog.Post{
				Title:  "Post Title",
				Author: "Author Name",
				Time:   toTime("2021-04-04T22:00:00Z"),
			})
	})

	t.Run("It parses content body after separator into Post Content", func(t *testing.T) {
		assertParsedContent(t, ""+
			"--\n"+
			"Content\n",
			blog.Post{
				Content: "" +
					"Content\n",
			})

		assertParsedContent(t, ""+
			"--\n"+
			"Only first separator is considered\n"+
			"--\n"+
			"After second separator\n",
			blog.Post{
				Content: "" +
					"Only first separator is considered\n" +
					"--\n" +
					"After second separator\n",
			})

	})

	t.Run("It parses everything together", func(t *testing.T) {
		assertParsedContent(t, ""+
			"title: Post Title\n"+
			"author: Author Name\n"+
			"time: 2021-04-04 22:00\n"+
			"--\n"+
			"## Subtitle\n"+
			"\n"+
			"Content\n"+
			"--\n"+
			"- list\n",
			blog.Post{
				Title:  "Post Title",
				Author: "Author Name",
				Time:   toTime("2021-04-04T22:00:00Z"),
				Content: "" +
					"## Subtitle\n" +
					"\n" +
					"Content\n" +
					"--\n" +
					"- list\n",
			})
	})

	t.Run("It ignores content without separator", func(t *testing.T) {
		assertParseError(t, "", filesystem.ErrInvalidFormat)
		assertParseError(t, "Content\n", filesystem.ErrInvalidFormat)
		assertParseError(t, "author: Author\n", filesystem.ErrInvalidFormat)
	})

	t.Run("It returns error if time is in an invalid format", func(t *testing.T) {
		assertParseError(t, "time: 04/04/2021\n--\n", filesystem.ErrInvalidTime)
	})
}

func assertParsedContent(t *testing.T, content string, expectedPost blog.Post) {
	t.Helper()
	post, err := filesystem.ParseFileContent(content)
	assert.DeepEqual(t, expectedPost, post)
	assert.Nil(t, err)
}

func assertParseError(t *testing.T, content string, expectedError error) {
	t.Helper()
	post, err := filesystem.ParseFileContent(content)
	assert.DeepEqual(t, blog.Post{}, post)
	assert.Equal(t, expectedError, err)
}

func toTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
