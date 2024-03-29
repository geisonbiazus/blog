package handlers_test

import (
	"errors"
	"net/http"
	"regexp"
	"testing"

	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/internal/ui/web/handlers"
	"github.com/geisonbiazus/blog/internal/ui/web/test"
	"github.com/geisonbiazus/blog/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

type feedHandlerFixture struct {
	usecase *listPostUseCaseSpy
	handler http.Handler
}

func TestFeedPostsHandler(t *testing.T) {
	setup := func() *feedHandlerFixture {
		usecase := &listPostUseCaseSpy{}
		templateRenderer := test.NewTestTemplateRenderer()
		baseURL := "http://example.com"
		handler := handlers.NewFeedHandler(usecase, templateRenderer, baseURL)

		return &feedHandlerFixture{
			usecase: usecase,
			handler: handler,
		}
	}

	t.Run("Given a list of posts exists it returns the feed containing all posts", func(t *testing.T) {
		f := setup()

		f.usecase.ReturnPosts = []blog.RenderedPost{renderedPost2, renderedPost1}

		res := test.DoGetRequest(f.handler, "/feed.atom")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/atom+xml", res.Header.Get("Content-Type"))
		assertFeedEqual(t, expectedFeed, body)
	})

	t.Run("Given no post exists it returns the empty feed", func(t *testing.T) {
		f := setup()

		f.usecase.ReturnPosts = []blog.RenderedPost{}

		res := test.DoGetRequest(f.handler, "/feed.atom")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/atom+xml", res.Header.Get("Content-Type"))
		assertFeedEqual(t, expectedEmptyFeed, body)
	})

	t.Run("Given an error occurs on getting the posts it returns 500", func(t *testing.T) {
		f := setup()

		f.usecase.ReturnError = errors.New("any error")

		res := test.DoGetRequest(f.handler, "/feed.atom")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Contains(t, body, "Internal server error")
	})
}

func assertFeedEqual(t *testing.T, expected, actual string) {
	t.Helper()
	assert.Equal(t, removeWhiteSpaces(expected), removeWhiteSpaces(actual))
}

var reWhiteSpace = regexp.MustCompile(`\s`)

func removeWhiteSpaces(source string) string {
	return reWhiteSpace.ReplaceAllString(source, "")
}

var expectedFeed = `<?xml version="1.0" encoding="UTF-8"?><feed xmlns="http://www.w3.org/2005/Atom">
	<title>Geison Biazus</title>
	<id>http://example.com</id>
	<updated>2021-04-04T14:33:00Z</updated>
	<subtitle>My personal blog about software development.</subtitle>
	<link href="http://example.com"></link>
	<author>
		<name>Geison Biazus</name>
		<email>geisonbiazus@gmail.com</email>
	</author>
	<entry>
		<title>Test Post 2</title>
		<updated>2021-04-04T14:33:00Z</updated>
		<id>tag:example.com,2021-04-04:/posts/test-post-2</id>
		<content type="html">Rendered content for post 2</content>
		<link href="http://example.com/posts/test-post-2" rel="alternate"></link>
		<summary type="html"></summary>
		<author>
			<name>Geison Biazus</name>
		</author>
	</entry>
	<entry>
		<title>Test Post 1</title>
		<updated>2021-04-05T18:47:00Z</updated>
		<id>tag:example.com,2021-04-05:/posts/test-post-1</id>
		<content type="html">Rendered content for post 1</content>
		<link href="http://example.com/posts/test-post-1" rel="alternate"></link>
		<summary type="html"></summary>
		<author>
			<name>Geison Biazus</name>
		</author>
	</entry>
</feed>`

var expectedEmptyFeed = `<?xml version="1.0" encoding="UTF-8"?><feed xmlns="http://www.w3.org/2005/Atom">
	<title>Geison Biazus</title>
	<id>http://example.com</id>
	<updated>2021-04-01T12:00:00Z</updated>
	<subtitle>My personal blog about software development.</subtitle>
	<link href="http://example.com"></link>
	<author>
		<name>Geison Biazus</name>
		<email>geisonbiazus@gmail.com</email>
	</author>
</feed>`
