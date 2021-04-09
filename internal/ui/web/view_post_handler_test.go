package web_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/internal/ui/web"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestViewPostHandler(t *testing.T) {
	t.Run("Given a valid post path it responds with the post HTML", func(t *testing.T) {
		usecase := &viewPostUseCaseSpy{}
		server := httptest.NewServer(web.NewRouter(usecase))
		defer server.Close()

		renderedPost := posts.RenderedPost{
			Title:   "post title",
			Author:  "post author",
			Time:    parseTime("2021-04-03T00:00:00+00:00"),
			Content: "<p>Content<p>",
		}

		usecase.RenderedPost = renderedPost

		res, _ := http.Get(server.URL + "/test-post")
		body, _ := io.ReadAll(res.Body)
		bodyString := string(body)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "test-post", usecase.ReceivedPath)

		assert.True(t, strings.Contains(bodyString, renderedPost.Title))
		assert.True(t, strings.Contains(bodyString, renderedPost.Author))
		assert.True(t, strings.Contains(bodyString, renderedPost.Content))

		fmt.Println(bodyString)
	})
}

func parseTime(timeString string) time.Time {
	t, _ := time.Parse(time.RFC3339, "2021-04-03T00:00:00+00:00")
	return t
}

type viewPostUseCaseSpy struct {
	ReceivedPath string
	RenderedPost posts.RenderedPost
	Error        error
}

func (u *viewPostUseCaseSpy) Run(path string) (posts.RenderedPost, error) {
	u.ReceivedPath = path
	return u.RenderedPost, u.Error
}
