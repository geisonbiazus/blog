package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/geisonbiazus/blog/internal/ui/web/lib"
	"github.com/geisonbiazus/blog/internal/ui/web/ports"
	"github.com/gorilla/feeds"
)

type FeedHandler struct {
	usecase  ports.ListPostUseCase
	template *lib.TemplateRenderer
	baseURL  string
}

func NewFeedHandler(usecase ports.ListPostUseCase, templateRenderer *lib.TemplateRenderer, baseURL string) *FeedHandler {
	return &FeedHandler{
		usecase:  usecase,
		template: templateRenderer,
		baseURL:  baseURL,
	}
}

func (h *FeedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	posts, err := h.usecase.Run()

	if err != nil {
		h.renderServerError(w)
	} else {
		h.renderFeed(w, posts)
	}
}

func (h *FeedHandler) renderServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	h.template.Render(w, "500.html", nil)
}

func (h *FeedHandler) renderFeed(w http.ResponseWriter, posts []blog.RenderedPost) {
	feed := h.buildFeed(posts)

	w.Header().Add("Content-Type", "application/atom+xml")
	w.WriteHeader(http.StatusOK)

	if err := feed.WriteAtom(w); err != nil {
		panic(fmt.Sprintf("Something went wrong rendering the feed: %v", err))
	}
}

func (h *FeedHandler) buildFeed(posts []blog.RenderedPost) *feeds.Feed {
	return &feeds.Feed{
		Title:       "Geison Biazus",
		Link:        &feeds.Link{Href: h.baseURL},
		Description: "My personal blog about software development.",
		Author:      &feeds.Author{Name: "Geison Biazus", Email: "geisonbiazus@gmail.com"},
		Created:     h.resolveUpdatedTime(posts),
		Items:       h.buildFeedItems(posts),
	}
}

func (h *FeedHandler) resolveUpdatedTime(posts []blog.RenderedPost) time.Time {
	if len(posts) == 0 {
		defaultTime, _ := time.Parse(time.RFC3339, "2021-04-01T12:00:00Z")
		return defaultTime
	}

	return posts[0].Post.Time
}

func (h *FeedHandler) buildFeedItems(posts []blog.RenderedPost) []*feeds.Item {
	items := []*feeds.Item{}

	for _, post := range posts {
		items = append(items, h.buildFeedItem(post))
	}

	return items
}

func (h *FeedHandler) buildFeedItem(post blog.RenderedPost) *feeds.Item {
	return &feeds.Item{
		Title:   post.Post.Title,
		Link:    &feeds.Link{Href: fmt.Sprintf("%s/posts/%s", h.baseURL, post.Post.Path)},
		Content: post.HTML,
		Author:  &feeds.Author{Name: post.Post.Author},
		Created: post.Post.Time,
	}
}
