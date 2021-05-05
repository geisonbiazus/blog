package web_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/internal/ui/web"
	"github.com/geisonbiazus/blog/pkg/assert"
	"github.com/geisonbiazus/blog/pkg/testhelper"
)

type feedHandlerFixture struct {
	usecase *listPostUseCaseSpy
	handler http.Handler
}

func TestFeedPostsHandler(t *testing.T) {
	setup := func() *feedHandlerFixture {
		usecase := &listPostUseCaseSpy{}
		templateRenderer := newTestTemplateRenderer()
		handler := web.NewFeedHandler(usecase, templateRenderer)

		return &feedHandlerFixture{
			usecase: usecase,
			handler: handler,
		}
	}

	t.Run("Given a list of posts exists it returns the feed containing all posts", func(t *testing.T) {
		f := setup()

		f.usecase.ReturnPosts = []blog.Post{post2, post1}

		res := doGetRequest(f.handler, "/feed.atom")
		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/atom+xml", res.Header.Get("Content-Type"))
		assert.Equal(t, cleanString(expectedFeed), cleanString(body))
	})
}

func cleanString(source string) string {
	source = strings.ReplaceAll(source, " ", "")
	source = strings.ReplaceAll(source, "\n", "")
	source = strings.ReplaceAll(source, "\t", "")
	return source
}

var expectedFeed = `<?xml version="1.0" encoding="UTF-8"?><feed xmlns="http://www.w3.org/2005/Atom">
	<title>Geison Biazus Blog</title>
	<id>https://blog.geisonbiazus.com</id>
	<updated>2021-04-04T14:33:00Z</updated>
	<subtitle>My personal blog about software development.</subtitle>
	<link href="https://blog.geisonbiazus.com"></link>
	<author>
		<name>Geison Biazus</name>
		<email>geisonbiazus@gmail.com</email>
	</author>
	<entry>
		<title>Test Post 2</title>
		<updated>2021-04-04T14:33:00Z</updated>
		<id>tag:blog.geisonbiazus.com,2021-04-04:/test-post-2</id>
		<content type="html">Content for post 2</content>
		<link href="https://blog.geisonbiazus.com/test-post-2" rel="alternate"></link>
		<summary type="html"></summary>
		<author>
			<name>Geison Biazus</name>
		</author>
	</entry>
	<entry>
		<title>Test Post 1</title>
		<updated>2021-04-05T18:47:00Z</updated>
		<id>tag:blog.geisonbiazus.com,2021-04-05:/test-post-1</id>
		<content type="html">Content for post 1</content>
		<link href="https://blog.geisonbiazus.com/test-post-1" rel="alternate"></link>
		<summary type="html"></summary>
		<author>
			<name>Geison Biazus</name>
		</author>
	</entry>
</feed>`
