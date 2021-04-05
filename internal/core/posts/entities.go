package posts

import "time"

type Post struct {
	Title   string
	Author  string
	Time    time.Time
	Path    string
	Content string
}

type RenderedPost struct {
	Title   string
	Author  string
	Time    time.Time
	Content string
}
