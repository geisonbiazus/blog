package discussion_test

import (
	"time"

	"github.com/geisonbiazus/blog/internal/core/discussion"
)

func newComment(
	params discussion.CommentParams, commentLoader discussion.CommentLoader,
) *discussion.Comment {
	return discussion.NewComment(discussion.CommentParams{
		ID:        stringOrDefault(params.ID, "ID"),
		SubjectID: stringOrDefault(params.SubjectID, "SUBJECT_ID"),
		Markdown:  stringOrDefault(params.Markdown, "Markdown"),
		HTML:      stringOrDefault(params.HTML, "HTML"),
		CreatedAt: timeOrDefault(params.CreatedAt, time.Now()),
	}, commentLoader)
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
