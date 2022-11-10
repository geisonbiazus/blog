package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/geisonbiazus/blog/pkg/dbrepo"
)

type CommentRepo struct {
	*dbrepo.Base
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{Base: dbrepo.NewBase(db)}
}

func (r *CommentRepo) SaveAuthor(ctx context.Context, author *discussion.Author) error {
	rows, err := r.Exec(ctx, `
		INSERT INTO discussion_authors 
			(id, name, avatar_url) 
		VALUES 
			($1, $2, $3)`,
		author.ID, author.Name, author.AvatarURL,
	)

	if err != nil {
		return fmt.Errorf("error on SaveAuthor when executing query: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("error on SaveAuthor, no affected rows")
	}

	return nil
}

func (r *CommentRepo) SaveComment(ctx context.Context, comment *discussion.Comment) error {
	rows, err := r.Exec(ctx, `
		INSERT INTO discussion_comments 
			(id, subject_id, author_id, markdown, html, created_at) 
		VALUES 
		($1, $2, $3, $4, $5, $6)`,
		comment.ID, comment.SubjectID, comment.AuthorID, comment.Markdown, comment.HTML, comment.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("error on SaveComment when executing query: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("error on SaveComment, no affected rows")
	}

	return nil
}

func (r *CommentRepo) GetCommentsAndRepliesRecursively(ctx context.Context, subjectID string) ([]*discussion.Comment, error) {
	return newGetCommentsAndRepliesRecursivelyQuery(r.Conn(ctx), ctx, subjectID).run()
}
