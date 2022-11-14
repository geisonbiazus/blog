package discussion

import (
	"time"
)

type Comment struct {
	ID        string
	SubjectID string
	AuthorID  string
	Author    *Author
	Markdown  string
	HTML      string
	CreatedAt time.Time
	Replies   []*Comment
}

func (c *Comment) Clone() *Comment {
	clone := *c
	return &clone
}

type Author struct {
	ID        string
	Name      string
	AvatarURL string
}
