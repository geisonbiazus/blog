package memory

import (
	"context"
	"sort"

	"github.com/geisonbiazus/blog/internal/core/discussion"
)

type CommentRepo struct {
	comments map[string]*discussion.Comment
}

func NewCommentRepo() *CommentRepo {
	return &CommentRepo{
		comments: make(map[string]*discussion.Comment),
	}
}

func (r *CommentRepo) Save(ctx context.Context, comment *discussion.Comment) error {
	r.comments[comment.ID] = comment
	return nil
}

func (r *CommentRepo) GetCommentsBySubjectID(ctx context.Context, subjectID string) ([]*discussion.Comment, error) {
	result := []*discussion.Comment{}

	for _, comment := range r.comments {
		if comment.SubjectID == subjectID {
			result = append(result, comment)
		}
	}

	sort.Sort(byCreatedAt(result))

	return result, nil
}

type byCreatedAt []*discussion.Comment

func (c byCreatedAt) Len() int           { return len(c) }
func (a byCreatedAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byCreatedAt) Less(i, j int) bool { return a[i].CreatedAt.Before(a[j].CreatedAt) }
