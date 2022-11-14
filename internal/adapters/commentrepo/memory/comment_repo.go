package memory

import (
	"context"
	"sort"

	"github.com/geisonbiazus/blog/internal/core/discussion"
)

type CommentRepo struct {
	comments map[string]*discussion.Comment
	authors  map[string]*discussion.Author
}

func NewCommentRepo() *CommentRepo {
	return &CommentRepo{
		comments: make(map[string]*discussion.Comment),
		authors:  make(map[string]*discussion.Author),
	}
}

func (r *CommentRepo) SaveComment(ctx context.Context, comment *discussion.Comment) error {
	r.comments[comment.ID] = comment
	return nil
}

func (r *CommentRepo) SaveAuthor(ctx context.Context, author *discussion.Author) error {
	r.authors[author.ID] = author
	return nil
}

func (r *CommentRepo) GetAuthorByID(ctx context.Context, id string) (*discussion.Author, error) {
	return r.authors[id], nil
}

func (r *CommentRepo) GetCommentsAndRepliesRecursively(ctx context.Context, subjectID string) ([]*discussion.Comment, error) {
	result := []*discussion.Comment{}

	for _, comment := range r.comments {
		if comment.SubjectID == subjectID {
			clone := comment.Clone()
			author, _ := r.GetAuthorByID(ctx, clone.AuthorID)
			clone.Author = author
			replies, _ := r.GetCommentsAndRepliesRecursively(ctx, comment.ID)
			clone.Replies = replies

			result = append(result, clone)
		}
	}

	sort.Sort(byCreatedAt(result))

	return result, nil
}

type byCreatedAt []*discussion.Comment

func (c byCreatedAt) Len() int           { return len(c) }
func (a byCreatedAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byCreatedAt) Less(i, j int) bool { return a[i].CreatedAt.Before(a[j].CreatedAt) }
