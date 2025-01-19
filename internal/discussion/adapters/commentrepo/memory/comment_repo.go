package memory

import (
	"context"
	"sort"

	"github.com/geisonbiazus/blog/internal/discussion/entities"
)

type CommentRepo struct {
	comments map[string]*entities.Comment
	authors  map[string]*entities.Author
}

func NewCommentRepo() *CommentRepo {
	return &CommentRepo{
		comments: make(map[string]*entities.Comment),
		authors:  make(map[string]*entities.Author),
	}
}

func (r *CommentRepo) SaveComment(ctx context.Context, comment *entities.Comment) error {
	r.comments[comment.ID] = comment
	return nil
}

func (r *CommentRepo) SaveAuthor(ctx context.Context, author *entities.Author) error {
	r.authors[author.ID] = author
	return nil
}

func (r *CommentRepo) GetAuthorByID(ctx context.Context, id string) (*entities.Author, error) {
	return r.authors[id], nil
}

func (r *CommentRepo) GetAuthorByUserID(ctx context.Context, userID string) (*entities.Author, error) {
	for _, author := range r.authors {
		if author.UserID == userID {
			return author, nil
		}
	}

	return nil, nil
}

func (r *CommentRepo) GetCommentsAndRepliesRecursively(ctx context.Context, subjectID string) ([]*entities.Comment, error) {
	result := []*entities.Comment{}

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

type byCreatedAt []*entities.Comment

func (c byCreatedAt) Len() int           { return len(c) }
func (a byCreatedAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byCreatedAt) Less(i, j int) bool { return a[i].CreatedAt.Before(a[j].CreatedAt) }
