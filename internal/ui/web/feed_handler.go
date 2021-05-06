package web

import (
	"fmt"
	"net/http"

	"github.com/geisonbiazus/blog/internal/core/blog"
	"github.com/gorilla/feeds"
)

type FeedHandler struct {
	usecase  ListPostUseCase
	template *TemplateRenderer
}

func NewFeedHandler(usecase ListPostUseCase, templateRenderer *TemplateRenderer) *FeedHandler {
	return &FeedHandler{
		usecase:  usecase,
		template: templateRenderer,
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
		Link:        &feeds.Link{Href: "https://blog.geisonbiazus.com"},
		Description: "My personal blog about software development.",
		Author:      &feeds.Author{Name: "Geison Biazus", Email: "geisonbiazus@gmail.com"},
		Created:     posts[0].Time,
		Items:       h.buildFeedItems(posts),
	}
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
		Link:    &feeds.Link{Href: fmt.Sprintf("https://blog.geisonbiazus.com/%s", post.Path)},
		Content: post.Content,
		Author:  &feeds.Author{Name: post.Author},
		Created: post.Time,
	}
}
