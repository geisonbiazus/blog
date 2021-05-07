package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/gorilla/feeds"
)

type FeedHandler struct {
	usecase  ListPostUseCase
	template *TemplateRenderer
	baseURL  string
}

func NewFeedHandler(usecase ListPostUseCase, templateRenderer *TemplateRenderer, baseURL string) *FeedHandler {
	return &FeedHandler{
		usecase:  usecase,
		template: templateRenderer,
		baseURL:  baseURL,
	}
}

func (h *FeedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	posts, _ := h.usecase.Run()
	feed := h.buildFeed(posts)

	w.Header().Add("Content-Type", "application/atom+xml")
	w.WriteHeader(http.StatusOK)
	feed.WriteAtom(w)
}

func (h *FeedHandler) buildFeed(posts []blog.Post) *feeds.Feed {
	return &feeds.Feed{
		Title:       "Geison Biazus Blog",
		Link:        &feeds.Link{Href: h.baseURL},
		Description: "My personal blog about software development.",
		Author:      &feeds.Author{Name: "Geison Biazus", Email: "geisonbiazus@gmail.com"},
		Created:     h.resolveUpdatedTime(posts),
		Items:       h.buildFeedItems(posts),
	}
}

func (h *FeedHandler) resolveUpdatedTime(posts []blog.Post) time.Time {
	if len(posts) == 0 {
		defaultTime, _ := time.Parse(time.RFC3339, "2021-04-01T12:00:00Z")
		return defaultTime
	}

	return posts[0].Time
}

func (h *FeedHandler) buildFeedItems(posts []blog.Post) []*feeds.Item {
	items := []*feeds.Item{}

	for _, post := range posts {
		items = append(items, h.buildFeedItem(post))
	}

	return items
}

func (h *FeedHandler) buildFeedItem(post blog.Post) *feeds.Item {
	return &feeds.Item{
		Title:   post.Title,
		Link:    &feeds.Link{Href: fmt.Sprintf("%s/%s", h.baseURL, post.Path)},
		Content: post.Content,
		Author:  &feeds.Author{Name: post.Author},
		Created: post.Time,
	}
}
