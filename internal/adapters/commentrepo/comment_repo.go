package commentrepo

import (
	"database/sql"

	"github.com/geisonbiazus/blog/internal/adapters/commentrepo/memory"
	"github.com/geisonbiazus/blog/internal/adapters/commentrepo/postgres"
)

func NewMemoryCommentRepo() *memory.CommentRepo {
	return memory.NewCommentRepo()
}

func NewPostgresCommentRepo(db *sql.DB) *postgres.CommentRepo {
	return postgres.NewCommentRepo(db)
}
