package web_test

import (
	"net/http"
	"testing"

	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/internal/ui/web"
	"github.com/geisonbiazus/blog/pkg/assert"
	"github.com/geisonbiazus/blog/pkg/testhelper"
)

type listPostHandlerFixture struct {
	usecase *listPostUseCaseSpy
	handler http.Handler
}

func TestListPostHandler(t *testing.T) {
	setup := func() *listPostHandlerFixture {
		usecase := &listPostUseCaseSpy{}
		templateRenderer := newTestTemplateRenderer()
		handler := web.NewListPostHandler(usecase, templateRenderer)

		return &listPostHandlerFixture{
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
		assert.Contains(t, body, post1.Title)
	})
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
