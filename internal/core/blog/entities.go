package blog

import "time"

type Post struct {
	Title    string
	Author   string
	Time     time.Time
	Path     string
	Markdown string
}

type RenderedPost struct {
	Post Post
	HTML string
}
