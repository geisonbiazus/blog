package web

import (
	"fmt"
	"net/http"

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

	feed := &feeds.Feed{
		Title:       "Geison Biazus Blog",
		Link:        &feeds.Link{Href: "https://blog.geisonbiazus.com"},
		Description: "My personal blog about software development.",
		Author:      &feeds.Author{Name: "Geison Biazus", Email: "geisonbiazus@gmail.com"},
		Created:     posts[0].Time,
	}

	for _, post := range posts {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:   post.Title,
			Link:    &feeds.Link{Href: fmt.Sprintf("https://blog.geisonbiazus.com/%s", post.Path)},
			Content: post.Content,
			Author:  &feeds.Author{Name: post.Author},
			Created: post.Time,
		})
	}

	w.Header().Add("Content-Type", "application/atom+xml")
	w.WriteHeader(http.StatusOK)
	feed.WriteAtom(w)
}
