package web_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/internal/ui/web/testhelper"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type viewPostHandlerFixture struct {
	usecase *testhelper.ViewPostUseCaseSpy
	server  *httptest.Server
}

func TestViewPostHandler(t *testing.T) {
	setup := func() *viewPostHandlerFixture {
		server, usecases := testhelper.NewTestServer()
		usecase := usecases.ViewPost.(*testhelper.ViewPostUseCaseSpy)

		return &viewPostHandlerFixture{
			usecase: usecase,
			server:  server,
		}
	}

	t.Run("Given a valid post path it responds with the post HTML", func(t *testing.T) {
		f := setup()
		defer f.server.Close()

		renderedPost := posts.RenderedPost{
			Title:   "post title",
			Author:  "post author",
			Time:    testhelper.ParseTime("2021-04-03T00:00:00+00:00"),
			Content: "<p>Content<p>",
		}

		f.usecase.RenderedPost = renderedPost
		res, _ := http.Get(f.server.URL + "/test-post")

		body := testhelper.ReadBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "test-post", f.usecase.ReceivedPath)

		assert.True(t, strings.Contains(body, renderedPost.Title))
		assert.True(t, strings.Contains(body, renderedPost.Author))
		assert.True(t, strings.Contains(body, renderedPost.Time.Format("02 Jan 06")))
		assert.True(t, strings.Contains(body, renderedPost.Content))
	})

	t.Run("Given a wrong post path it responds with not found", func(t *testing.T) {
		f := setup()
		defer f.server.Close()

		f.usecase.Error = posts.ErrPostNotFound

		res, _ := http.Get(f.server.URL + "/wrong-path")

		body := testhelper.ReadBody(res)

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Equal(t, "wrong-path", f.usecase.ReceivedPath)

		assert.True(t, strings.Contains(body, "Page not found"))
	})

	t.Run("Returns server error when sother error happens", func(t *testing.T) {
		f := setup()
		defer f.server.Close()

		f.usecase.Error = errors.New("any error")

		res, _ := http.Get(f.server.URL + "/post-path")

		body := testhelper.ReadBody(res)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, "post-path", f.usecase.ReceivedPath)

		assert.True(t, strings.Contains(body, "Internal server error"))
	})
}
