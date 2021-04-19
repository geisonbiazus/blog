package web_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/internal/ui/web"
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
		templateRenderer := newTestTemplateRenderer()
		handler := web.NewListPostsHandler(usecase, templateRenderer)

		return &listPostsHandlerFixture{
			usecase: usecase,
			handler: handler,
		}
	}

	t.Run("Given a list of posts exists it renders the posts", func(t *testing.T) {
		f := setup()

		f.usecase.ReturnPosts = []posts.Post{post2, post1}

		res := doGetRequest(f.handler, "/index")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assertContainsListedPost(t, body, post2)
		assertContainsListedPost(t, body, post1)
	})

	t.Run("It renders server error when and error is returned", func(t *testing.T) {
		f := setup()

		f.usecase.ReturnError = errors.New("some error")

		res := doGetRequest(f.handler, "/index")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Contains(t, body, "Internal server error")
	})
}

func assertContainsListedPost(t *testing.T, body string, post posts.Post) {
	assert.Contains(t, body, post.Title)
	assert.Contains(t, body, post.Author)
	assert.Contains(t, body, post.Time.Format("02 Jan 06"))
	assert.Contains(t, body, post.Path)
}

type listPostUseCaseSpy struct {
	ReturnPosts []posts.Post
	ReturnError error
}

func (u *listPostUseCaseSpy) Run() ([]posts.Post, error) {
	return u.ReturnPosts, u.ReturnError
}

var post1 = posts.Post{
	Title:  "Test Post 1",
	Author: "Geison Biazus",
	Path:   "test-post-1",
	Time:   testhelper.ParseTime("2021-04-05T18:47:00Z"),
}

var post2 = posts.Post{
	Title:  "Test Post 2",
	Author: "Geison Biazus",
	Path:   "test-post-2",
	Time:   testhelper.ParseTime("2021-04-04T14:33:00Z"),
}
