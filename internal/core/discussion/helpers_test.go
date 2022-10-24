package discussion_test

import (
	"time"

	"github.com/geisonbiazus/blog/internal/core/discussion"
)

func newComment(params discussion.Comment) *discussion.Comment {
	return &discussion.Comment{
		ID:        stringOrDefault(params.ID, "ID"),
		SubjectID: stringOrDefault(params.SubjectID, "SUBJECT_ID"),
		AuthorID:  stringOrDefault(params.AuthorID, "AUTHOR_ID"),
		Markdown:  stringOrDefault(params.Markdown, "Markdown"),
		HTML:      stringOrDefault(params.HTML, "HTML"),
		CreatedAt: timeOrDefault(params.CreatedAt, time.Now()),
		Replies:   sliceOrDefault(params.Replies, []*discussion.Comment{}),
	}
}

func stringOrDefault(value, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

func timeOrDefault(value, defaultValue time.Time) time.Time {
	if (value != time.Time{}) {
		return value
	}
	return defaultValue
}

func sliceOrDefault[T any](value []T, defaultValue []T) []T {
	if value != nil {
		return value
	}

	return defaultValue
}
