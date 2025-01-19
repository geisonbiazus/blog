package test

import (
	"time"

	"github.com/geisonbiazus/blog/internal/discussion/entities"
)

func NewComment(params entities.Comment) *entities.Comment {
	return &entities.Comment{
		ID:        valueOrDefault(params.ID, "COMMENT_ID"),
		SubjectID: valueOrDefault(params.SubjectID, "SUBJECT_ID"),
		AuthorID:  valueOrDefault(params.AuthorID, "AUTHOR_ID"),
		Author:    valueOrDefault(params.Author, nil),
		Markdown:  valueOrDefault(params.Markdown, "Markdown"),
		HTML:      valueOrDefault(params.HTML, "HTML"),
		CreatedAt: valueOrDefault(params.CreatedAt, time.Now()),
		Replies:   sliceOrDefault(params.Replies, []*entities.Comment{}),
	}
}

func NewAuthor(params entities.Author) *entities.Author {
	return &entities.Author{
		ID:        valueOrDefault(params.ID, "AUTHOR_ID"),
		UserID:    valueOrDefault(params.UserID, "USER_ID"),
		Name:      valueOrDefault(params.Name, "Author"),
		AvatarURL: valueOrDefault(params.AvatarURL, "http://example.com/avatar"),
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
