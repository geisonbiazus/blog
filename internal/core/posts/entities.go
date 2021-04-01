package posts

import "time"

type Post struct {
	Title   string
	Authors []string
	Time    time.Time
	Path    string
	Content string
}

type RenderedPost struct {
	Title   string
	Authors []string
	Time    time.Time
	Content string
}
