package commentrepo

import "github.com/geisonbiazus/blog/internal/adapters/commentrepo/memory"

func NewMemoryCommentRepo() *memory.CommentRepo {
	return memory.NewCommentRepo()
}
