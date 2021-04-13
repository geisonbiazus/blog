package web_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/internal/ui/web"
	"github.com/geisonbiazus/blog/pkg/assert"
	"github.com/geisonbiazus/blog/pkg/testhelper"
)

type viewPostHandlerFixture struct {
	usecase *viewPostUseCaseSpy
	handler http.Handler
}

func TestViewPostHandler(t *testing.T) {
	setup := func() *viewPostHandlerFixture {
		usecase := &viewPostUseCaseSpy{}
		templateRenderer := newTestTemplateRenderer()
		handler := web.NewViewPostHandler(usecase, templateRenderer)

		return &viewPostHandlerFixture{
			usecase: usecase,
			handler: handler,
		}
	}

	t.Run("Given a valid post path it responds with the post HTML", func(t *testing.T) {
		f := setup()

		renderedPost := buildRenderedPost()
		f.usecase.RenderedPost = renderedPost

		res := doGetRequest(f.handler, "/post-path")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "post-path", f.usecase.ReceivedPath)
		assertContainsRenderedPost(t, body, renderedPost)
	})

	t.Run("Given a wrong post path it responds with not found", func(t *testing.T) {
		f := setup()

		f.usecase.Error = posts.ErrPostNotFound

		res := doGetRequest(f.handler, "/post-path")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Equal(t, "post-path", f.usecase.ReceivedPath)
		assert.True(t, strings.Contains(body, "Page not found"))
	})

	t.Run("Returns server error when other error happens", func(t *testing.T) {
		f := setup()

		f.usecase.Error = errors.New("any error")

		res := doGetRequest(f.handler, "/post-path")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, "post-path", f.usecase.ReceivedPath)
		assert.True(t, strings.Contains(body, "Internal server error"))
	})
}

func buildRenderedPost() posts.RenderedPost {
	return posts.RenderedPost{
		Title:   "post title",
		Author:  "post author",
		Time:    testhelper.ParseTime("2021-04-03T00:00:00+00:00"),
		Content: "<p>Content<p>",
	}
}

func doGetRequest(handler http.Handler, path string) *http.Response {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rw := httptest.NewRecorder()

	handler.ServeHTTP(rw, req)

	return rw.Result()
}

func assertContainsRenderedPost(t *testing.T, body string, renderedPost posts.RenderedPost) {
	assert.True(t, strings.Contains(body, renderedPost.Title))
	assert.True(t, strings.Contains(body, renderedPost.Author))
	assert.True(t, strings.Contains(body, renderedPost.Time.Format("02 Jan 06")))
	assert.True(t, strings.Contains(body, renderedPost.Content))
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
