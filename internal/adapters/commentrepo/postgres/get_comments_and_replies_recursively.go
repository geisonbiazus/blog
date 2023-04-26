package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/geisonbiazus/blog/internal/core/discussion"
	"github.com/geisonbiazus/blog/pkg/dbrepo"
)

type getCommentsAndRepliesRecursivelyQuery struct {
	conn       dbrepo.Connection
	ctx        context.Context
	subjectID  string
	result     []*discussion.Comment
	rows       *sql.Rows
	commentMap map[string][]*discussion.Comment
}

func newGetCommentsAndRepliesRecursivelyQuery(conn dbrepo.Connection, ctx context.Context, subjectID string) *getCommentsAndRepliesRecursivelyQuery {
	return &getCommentsAndRepliesRecursivelyQuery{
		conn:      conn,
		ctx:       ctx,
		subjectID: subjectID,
	}
}

func (q *getCommentsAndRepliesRecursivelyQuery) run() ([]*discussion.Comment, error) {
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
	q.result = []*discussion.Comment{}
	q.commentMap = map[string][]*discussion.Comment{}
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

func (q *getCommentsAndRepliesRecursivelyQuery) scanRow(row *sql.Rows) (*discussion.Comment, error) {
	comment := &discussion.Comment{
		Author: &discussion.Author{Persisted: true},
	}

	err := q.rows.Scan(
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

func (q *getCommentsAndRepliesRecursivelyQuery) addToCommentMap(comment *discussion.Comment) {
	if q.commentMap[comment.SubjectID] == nil {
		q.commentMap[comment.SubjectID] = []*discussion.Comment{}
	}

	q.commentMap[comment.SubjectID] = append(q.commentMap[comment.SubjectID], comment)
}

func (q *getCommentsAndRepliesRecursivelyQuery) buildResult() {
	q.result = q.commentMap[q.subjectID]
	q.appendRepliesRecursively(q.result)
}

func (q *getCommentsAndRepliesRecursivelyQuery) appendRepliesRecursively(comments []*discussion.Comment) {
	for _, comment := range comments {
		comment.Replies = q.commentMap[comment.ID]

		if comment.Replies == nil {
			comment.Replies = []*discussion.Comment{}
		}

		q.appendRepliesRecursively(comment.Replies)
	}
}
