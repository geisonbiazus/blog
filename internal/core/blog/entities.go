package blog

import (
	"errors"
	"time"
)

type Post struct {
	Title       string
	Author      string
	Time        time.Time
	Path        string
	Description string
	ImagePath   string
	Markdown    string
}

type RenderedPost struct {
	Post Post
	HTML string
}

var ErrPostNotFound = errors.New("post not found")
