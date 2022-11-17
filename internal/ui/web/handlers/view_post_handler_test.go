package handlers_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/geisonbiazus/blog/internal/ui/web/handlers"
	"github.com/geisonbiazus/blog/internal/ui/web/lib"
	"github.com/geisonbiazus/blog/internal/ui/web/test"
	"github.com/geisonbiazus/blog/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

type viewPostHandlerFixture struct {
	viewPostUseCase     *viewPostUseCaseSpy
	listCommentsUseCase *listCommentsUseCaseSpy
	handler             http.Handler
}

func TestViewPostHandler(t *testing.T) {
	setup := func() *viewPostHandlerFixture {
		viewPostUseCase := &viewPostUseCaseSpy{}
		listCommentsUseCase := &listCommentsUseCaseSpy{}
		templateRenderer := test.NewTestTemplateRenderer()
		handler := handlers.NewViewPostHandler(viewPostUseCase, listCommentsUseCase, templateRenderer)

		return &viewPostHandlerFixture{
			viewPostUseCase:     viewPostUseCase,
			listCommentsUseCase: listCommentsUseCase,
			handler:             handler,
		}
	}

	t.Run("Given a valid post path it responds with the post HTML", func(t *testing.T) {
		f := setup()

		renderedPost := buildRenderedPost()
		f.viewPostUseCase.ReturnPost = renderedPost

		res := test.DoGetRequest(f.handler, "/posts/post-path")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "post-path", f.viewPostUseCase.ReceivedPath)
		assertContainsRenderedPost(t, body, renderedPost)
	})

	t.Run("Given a post with comments it renders the comments and replies", func(t *testing.T) {
		f := setup()

		renderedPost := buildRenderedPost()
		comments := buildComments()
		f.viewPostUseCase.ReturnPost = renderedPost
		f.listCommentsUseCase.ReturnComments = comments

		res := test.DoGetRequest(f.handler, "/posts/post-path")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assertContainsComments(t, body, comments)
	})

	t.Run("Given a post with no comments it doesn't render comments", func(t *testing.T) {
		f := setup()

		renderedPost := buildRenderedPost()
		f.viewPostUseCase.ReturnPost = renderedPost
		f.listCommentsUseCase.ReturnComments = []*discussion.Comment{}

		res := test.DoGetRequest(f.handler, "/posts/post-path")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.NotContains(t, body, "Comments")
	})

	t.Run("Given an error is returned when loading coments it responds with server error", func(t *testing.T) {
		f := setup()

		renderedPost := buildRenderedPost()
		f.viewPostUseCase.ReturnPost = renderedPost
		f.listCommentsUseCase.ReturnError = errors.New("any error")

		res := test.DoGetRequest(f.handler, "/posts/post-path")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.True(t, strings.Contains(body, "Internal server error"))
	})

	t.Run("Given a wrong post path it responds with not found", func(t *testing.T) {
		f := setup()

		f.viewPostUseCase.ReturnError = blog.ErrPostNotFound

		res := test.DoGetRequest(f.handler, "/posts/post-path")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Equal(t, "post-path", f.viewPostUseCase.ReceivedPath)
		assert.True(t, strings.Contains(body, "Page not found"))
	})

	t.Run("Returns server error when other error happens", func(t *testing.T) {
		f := setup()

		f.viewPostUseCase.ReturnError = errors.New("any error")

		res := test.DoGetRequest(f.handler, "/posts/post-path")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, "post-path", f.viewPostUseCase.ReceivedPath)
		assert.True(t, strings.Contains(body, "Internal server error"))
	})
}

func buildRenderedPost() blog.RenderedPost {
	return blog.RenderedPost{
		Post: blog.Post{
			Title:       "post title",
			Author:      "post author",
			Path:        "post-path",
			Description: "post description",
			ImagePath:   "/static/image/post.png",
			Time:        testhelper.ParseTime("2021-04-03T00:00:00+00:00"),
		},
		HTML: "<p>Content<p>",
	}
}

func buildComments() []*discussion.Comment {
	return []*discussion.Comment{
		{
			ID:        "COMMENT_ID",
			SubjectID: "post-path",
			Author: &discussion.Author{
				ID:        "COMMENT_AUTHOR_ID",
				Name:      "Comment Author",
				AvatarURL: "https://example.com/comment-author-avatar",
			},
			HTML:      "Comment HTML",
			CreatedAt: testhelper.ParseTime("2021-04-04T00:00:00+00:00"),
			Replies: []*discussion.Comment{
				{
					ID:        "REPLY_ID",
					SubjectID: "COMMENT_ID",
					Author: &discussion.Author{
						ID:        "REPLY_AUTHOR_ID",
						Name:      "Reply Author",
						AvatarURL: "https://example.com/reply-author-avatar",
					},
					HTML:      "Comment HTML",
					CreatedAt: testhelper.ParseTime("2021-04-05T00:00:00+00:00"),
				},
			},
		},
	}
}

func assertContainsRenderedPost(t *testing.T, body string, renderedPost blog.RenderedPost) {
	assert.Contains(t, body, renderedPost.Post.Title)
	assert.Contains(t, body, renderedPost.Post.Author)
	assert.Contains(t, body, renderedPost.Post.Description)
	assert.Contains(t, body, fmt.Sprintf("http://example.com%s", renderedPost.Post.ImagePath))
	assert.Contains(t, body, fmt.Sprintf("http://example.com/posts/%s", renderedPost.Post.Path))
	assert.Contains(t, body, renderedPost.Post.Time.Format(lib.DateFormat))
	assert.Contains(t, body, renderedPost.HTML)
}

func assertContainsComments(t *testing.T, body string, comments []*discussion.Comment) {
	assert.Contains(t, body, "Comments")
	for _, comment := range comments {
		assert.Contains(t, body, comment.Author.Name)
		assert.Contains(t, body, comment.Author.AvatarURL)
		assert.Contains(t, body, comment.CreatedAt.Format(lib.DateFormat))
		assert.Contains(t, body, comment.HTML)
		if comment.Replies != nil {
			assertContainsComments(t, body, comment.Replies)
		}
	}
}

type viewPostUseCaseSpy struct {
	ReceivedPath string
	ReturnPost   blog.RenderedPost
	ReturnError  error
}

func (u *viewPostUseCaseSpy) Run(path string) (blog.RenderedPost, error) {
	u.ReceivedPath = path
	return u.ReturnPost, u.ReturnError
}

type listCommentsUseCaseSpy struct {
	ReceivedCtx       context.Context
	ReceivedSubjectID string
	ReturnComments    []*discussion.Comment
	ReturnError       error
}

func (u *listCommentsUseCaseSpy) Run(ctx context.Context, subjectID string) ([]*discussion.Comment, error) {
	u.ReceivedCtx = ctx
	u.ReceivedSubjectID = subjectID
	return u.ReturnComments, u.ReturnError
}
