package discussion_test

import (
	"time"

	"github.com/geisonbiazus/blog/internal/core/discussion"
)

func newComment(params discussion.Comment) *discussion.Comment {
	return &discussion.Comment{
		ID:        valueOrDefault(params.ID, "ID"),
		SubjectID: valueOrDefault(params.SubjectID, "SUBJECT_ID"),
		AuthorID:  valueOrDefault(params.AuthorID, "AUTHOR_ID"),
		Author:    valueOrDefault(params.Author, nil),
		Markdown:  valueOrDefault(params.Markdown, "Markdown"),
		HTML:      valueOrDefault(params.HTML, "HTML"),
		CreatedAt: valueOrDefault(params.CreatedAt, time.Now()),
		Replies:   sliceOrDefault(params.Replies, []*discussion.Comment{}),
	}
}

func valueOrDefault[T comparable](value T, defaultValue T) T {
	var empty T
	if value != empty {
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
