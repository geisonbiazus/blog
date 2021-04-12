package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/postrepo/filesystem"
	"github.com/geisonbiazus/blog/internal/adapters/renderer/goldmark"
	"github.com/geisonbiazus/blog/internal/core/posts"
	"github.com/geisonbiazus/blog/internal/ui/web"
	"github.com/geisonbiazus/blog/pkg/assert"
	"github.com/geisonbiazus/blog/pkg/testhelper"
)

func TestViewPostIntegration(t *testing.T) {
	setup := func() *httptest.Server {
		templateRenderer, err := web.NewTemplateRenderer("../../web/template")

		if err != nil {
			t.Fatal(err)
		}

		postRepo := filesystem.NewPostRepo("../posts")
		renderer := goldmark.NewGoldmarkRenderer()
		viewPostUseCase := posts.NewVewPostUseCase(postRepo, renderer)

		usecases := &web.UseCases{
			ViewPost: viewPostUseCase,
		}

		router := web.NewRouter(templateRenderer, usecases)

		server := httptest.NewServer(router)

		return server
	}

	t.Run("Given a valid post path it responds with the post HTML", func(t *testing.T) {
		server := setup()
		defer server.Close()

		res, _ := http.Get(server.URL + "/test-post")

		body := testhelper.ReadResponseBody(res)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.True(t, strings.Contains(body, "Test Post"))
		assert.True(t, strings.Contains(body, "Geison Biazus"))
		assert.True(t, strings.Contains(body, "05 Apr 21"))
		assert.True(t, strings.Contains(body, "Content"))
	})
}
