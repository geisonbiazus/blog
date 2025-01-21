package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/geisonbiazus/blog/internal/discussion/entities"
	"github.com/geisonbiazus/blog/pkg/dbrepo"
)

type getCommentsAndRepliesRecursivelyQuery struct {
	conn       dbrepo.Connection
	ctx        context.Context
	subjectID  string
	result     []*entities.Comment
	rows       *sql.Rows
	commentMap map[string][]*entities.Comment
}

func newGetCommentsAndRepliesRecursivelyQuery(conn dbrepo.Connection, ctx context.Context, subjectID string) *getCommentsAndRepliesRecursivelyQuery {
	return &getCommentsAndRepliesRecursivelyQuery{
		conn:      conn,
		ctx:       ctx,
		subjectID: subjectID,
	}
}

func (q *getCommentsAndRepliesRecursivelyQuery) run() ([]*entities.Comment, error) {
	q.initializeVariables()

	if err := q.executeQuery(); err != nil {
		return q.result, err
	}

	if err := q.scanRowsAndBuildCommentMap(); err != nil {
		return q.result, err
	}

	q.buildResult()

	return q.result, nil
}

func (q *getCommentsAndRepliesRecursivelyQuery) initializeVariables() {
	q.result = []*entities.Comment{}
	q.commentMap = map[string][]*entities.Comment{}
	q.rows = nil
}

func (q *getCommentsAndRepliesRecursivelyQuery) executeQuery() error {
	rows, err := q.conn.QueryContext(q.ctx, `
		WITH RECURSIVE comments_and_replies as (
			SELECT 
				c.id, c.subject_id, c.author_id, c.markdown, c.html, c.created_at,
				a.id AS author_id, a.auth_user_id, a.name, a.avatar_url
			FROM discussion_comments c
			JOIN discussion_authors a ON c.author_id = a.id
			WHERE c.subject_id = $1
			
			UNION

			SELECT 
				c.id, c.subject_id, c.author_id, c.markdown, c.html, c.created_at,
				a.id AS author_id, a.auth_user_id, a.name, a.avatar_url
			FROM discussion_comments c
			JOIN discussion_authors a ON c.author_id = a.id
			JOIN comments_and_replies cr ON c.subject_id = cr.id::TEXT
		) 
		SELECT * 
		FROM comments_and_replies
		ORDER BY created_at`,
		q.subjectID,
	)

	q.rows = rows

	if err != nil {
		return fmt.Errorf("error on GetCommentsAndRepliesRecursively when resolving query: %w", err)
	}

	return nil
}

func (q *getCommentsAndRepliesRecursivelyQuery) scanRowsAndBuildCommentMap() error {
	for q.rows.Next() {
		comment, err := q.scanRow(q.rows)
		if err != nil {
			return err
		}

		q.addToCommentMap(comment)
	}

	return nil
}

func (q *getCommentsAndRepliesRecursivelyQuery) scanRow(row *sql.Rows) (*entities.Comment, error) {
	comment := &entities.Comment{
		Author: &entities.Author{Persisted: true},
	}

	err := row.Scan(
		&comment.ID,
		&comment.SubjectID,
		&comment.AuthorID,
		&comment.Markdown,
		&comment.HTML,
		&comment.CreatedAt,
		&comment.Author.ID,
		&comment.Author.UserID,
		&comment.Author.Name,
		&comment.Author.AvatarURL,
	)

	if err != nil {
		return comment, fmt.Errorf("error on GetCommentsAndRepliesRecursively when scanning row: %w", err)
	}

	return comment, err
}

func (q *getCommentsAndRepliesRecursivelyQuery) addToCommentMap(comment *entities.Comment) {
	if q.commentMap[comment.SubjectID] == nil {
		q.commentMap[comment.SubjectID] = []*entities.Comment{}
	}

	q.commentMap[comment.SubjectID] = append(q.commentMap[comment.SubjectID], comment)
}

func (q *getCommentsAndRepliesRecursivelyQuery) buildResult() {
	q.result = q.commentMap[q.subjectID]
	q.appendRepliesRecursively(q.result)
}

func (q *getCommentsAndRepliesRecursivelyQuery) appendRepliesRecursively(comments []*entities.Comment) {
	for _, comment := range comments {
		comment.Replies = q.commentMap[comment.ID]

		if comment.Replies == nil {
			comment.Replies = []*entities.Comment{}
		}

		q.appendRepliesRecursively(comment.Replies)
	}
}
