package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/geisonbiazus/blog/internal/discussion/entities"
	"github.com/geisonbiazus/blog/pkg/dbrepo"
)

type CommentRepo struct {
	*dbrepo.Base
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{Base: dbrepo.NewBase(db)}
}

func (r *CommentRepo) SaveAuthor(ctx context.Context, author *entities.Author) error {
	if author.Persisted {
		return r.updateAuthor(ctx, author)
	} else {
		return r.insertAuthor(ctx, author)
	}
}

func (r *CommentRepo) insertAuthor(ctx context.Context, author *entities.Author) error {
	err := r.Insert(ctx, "discussion_authors", map[string]interface{}{
		"id":           author.ID,
		"auth_user_id": author.UserID,
		"name":         author.Name,
		"avatar_url":   author.AvatarURL,
	})

	if err != nil {
		return fmt.Errorf("error on insertAuthor: %w", err)
	}

	author.Persisted = true

	return nil
}

func (r *CommentRepo) updateAuthor(ctx context.Context, author *entities.Author) error {
	err := r.Update(ctx, "discussion_authors", author.ID, map[string]interface{}{
		"auth_user_id": author.UserID,
		"name":         author.Name,
		"avatar_url":   author.AvatarURL,
	})

	if err != nil {
		return fmt.Errorf("error on updateAuthor: %w", err)
	}

	return nil
}

func (r *CommentRepo) GetAuthorByID(ctx context.Context, id string) (*entities.Author, error) {
	conn := r.Conn(ctx)

	row := conn.QueryRowContext(ctx, `
		SELECT 
		id, auth_user_id, name, avatar_url 
		FROM discussion_authors 
		WHERE id = $1`,
		id,
	)

	author := &entities.Author{Persisted: true}

	err := row.Scan(&author.ID, &author.UserID, &author.Name, &author.AvatarURL)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error on GetAuthorByID when executing query: %w", err)
	}

	return author, nil
}

func (r *CommentRepo) GetAuthorByUserID(ctx context.Context, userID string) (*entities.Author, error) {
	conn := r.Conn(ctx)

	row := conn.QueryRowContext(ctx, `
		SELECT 
			id, auth_user_id, name, avatar_url 
		FROM discussion_authors 
		WHERE auth_user_id = $1`,
		userID,
	)

	author := &entities.Author{Persisted: true}

	err := row.Scan(&author.ID, &author.UserID, &author.Name, &author.AvatarURL)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error on GetAuthorByUserID when executing query: %w", err)
	}

	return author, nil
}

func (r *CommentRepo) SaveComment(ctx context.Context, comment *entities.Comment) error {
	err := r.Insert(ctx, "discussion_comments", map[string]interface{}{
		"id":         comment.ID,
		"subject_id": comment.SubjectID,
		"author_id":  comment.AuthorID,
		"markdown":   comment.Markdown,
		"html":       comment.HTML,
		"created_at": comment.CreatedAt,
	})

	if err != nil {
		return fmt.Errorf("error on SaveComment: %w", err)
	}

	return nil
}

func (r *CommentRepo) GetCommentsAndRepliesRecursively(ctx context.Context, subjectID string) ([]*entities.Comment, error) {
	return newGetCommentsAndRepliesRecursivelyQuery(r.Conn(ctx), ctx, subjectID).run()
}
