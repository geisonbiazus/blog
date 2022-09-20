package handlers_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/internal/ui/web/handlers"
	"github.com/geisonbiazus/blog/internal/ui/web/lib"
	"github.com/geisonbiazus/blog/internal/ui/web/test"
	"github.com/geisonbiazus/blog/pkg/assert"
	"github.com/geisonbiazus/blog/pkg/testhelper"
)

type listPostsHandlerFixture struct {
	usecase *listPostUseCaseSpy
	handler http.Handler
}

func TestListPostsHandler(t *testing.T) {
	setup := func() *listPostsHandlerFixture {
		usecase := &listPostUseCaseSpy{}
		templateRenderer := test.NewTestTemplateRenderer()
		handler := handlers.NewListPostsHandler(usecase, templateRenderer)

		return &listPostsHandlerFixture{
			usecase: usecase,
			handler: handler,
		}
	}

	t.Run("Given a list of posts exists it renders the posts", func(t *testing.T) {
		f := setup()

		f.usecase.ReturnPosts = []blog.RenderedPost{renderedPost1, renderedPost2}

		res := test.DoGetRequest(f.handler, "/index")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assertContainsListedPost(t, body, renderedPost2)
		assertContainsListedPost(t, body, renderedPost1)
	})

	t.Run("It renders server error when and error is returned", func(t *testing.T) {
		f := setup()

		f.usecase.ReturnError = errors.New("some error")

		res := test.DoGetRequest(f.handler, "/posts")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Contains(t, body, "Internal server error")
	})
}

func assertContainsListedPost(t *testing.T, body string, post blog.RenderedPost) {
	t.Helper()
	assert.Contains(t, body, post.Post.Title)
	assert.Contains(t, body, post.Post.Author)
	assert.Contains(t, body, post.Post.Time.Format(lib.DateFormat))
	assert.Contains(t, body, fmt.Sprintf("/posts/%s", post.Post.Path))
}

type listPostUseCaseSpy struct {
	ReturnPosts []blog.RenderedPost
	ReturnError error
}

func (u *listPostUseCaseSpy) Run() ([]blog.RenderedPost, error) {
	return u.ReturnPosts, u.ReturnError
}

var post1 = blog.Post{
	Title:    "Test Post 1",
	Author:   "Geison Biazus",
	Path:     "test-post-1",
	Time:     testhelper.ParseTime("2021-04-05T18:47:00Z"),
	Markdown: "Content for post 1",
}

var post2 = blog.Post{
	Title:    "Test Post 2",
	Author:   "Geison Biazus",
	Path:     "test-post-2",
	Time:     testhelper.ParseTime("2021-04-04T14:33:00Z"),
	Markdown: "Content for post 2",
}

var renderedPost1 = blog.RenderedPost{Post: post1, HTML: "Rendered content for post 1"}
var renderedPost2 = blog.RenderedPost{Post: post2, HTML: "Rendered content for post 2"}
