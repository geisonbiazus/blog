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
	conn := r.Conn(ctx)
	result := []*discussion.Comment{}

	rows, err := conn.QueryContext(ctx, `
		WITH RECURSIVE comments_and_replies as (
			SELECT 
				c.id, c.subject_id, c.author_id, c.markdown, c.html, c.created_at,
				a.id AS author_id, a.name, a.avatar_url
			FROM discussion_comments c
			JOIN discussion_authors a ON c.author_id = a.id
			WHERE c.subject_id = $1
			
			UNION

			SELECT 
				c.id, c.subject_id, c.author_id, c.markdown, c.html, c.created_at,
				a.id AS author_id, a.name, a.avatar_url
			FROM discussion_comments c
			JOIN discussion_authors a ON c.author_id = a.id
			JOIN comments_and_replies cr ON c.subject_id = cr.id::TEXT
		) 
		SELECT * 
		FROM comments_and_replies
		ORDER BY created_at`,
		subjectID,
	)

	if err != nil {
		return result, fmt.Errorf("error on GetCommentsAndRepliesRecursively when resolving query: %w", err)
	}

	commentMap := map[string][]*discussion.Comment{}

	for rows.Next() {
		comment := &discussion.Comment{
			Author: &discussion.Author{},
		}

		err = rows.Scan(
			&comment.ID,
			&comment.SubjectID,
			&comment.AuthorID,
			&comment.Markdown,
			&comment.HTML,
			&comment.CreatedAt,
			&comment.Author.ID,
			&comment.Author.Name,
			&comment.Author.AvatarURL,
		)

		if err != nil {
			return result, fmt.Errorf("error on GetCommentsAndRepliesRecursively when scanning row: %w", err)
		}

		if commentMap[comment.SubjectID] == nil {
			commentMap[comment.SubjectID] = []*discussion.Comment{}
		}

		commentMap[comment.SubjectID] = append(commentMap[comment.SubjectID], comment)

		if comment.SubjectID == subjectID {
			result = append(result, comment)
		}
	}

	r.appendRepliesRecursively(commentMap, result)

	return result, err
}

func (r *CommentRepo) appendRepliesRecursively(commentMap map[string][]*discussion.Comment, comments []*discussion.Comment) {
	for _, comment := range comments {
		comment.Replies = commentMap[comment.ID]

		if comment.Replies == nil {
			comment.Replies = []*discussion.Comment{}
		}

		r.appendRepliesRecursively(commentMap, comment.Replies)
	}
}
