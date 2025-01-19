package discussion

import (
	"github.com/geisonbiazus/blog/internal/app/shared"
	"github.com/geisonbiazus/blog/internal/discussion/adapters/commentrepo"
	"github.com/geisonbiazus/blog/internal/discussion/entities"
	"github.com/geisonbiazus/blog/internal/discussion/ports"
	"github.com/geisonbiazus/blog/internal/discussion/usecases"
)

type Comment = entities.Comment
type Author = entities.Author

type CommentRepo = ports.CommentRepo

type ListCommentsUseCase = usecases.ListCommentsUseCase
type SaveAuthorUseCase = usecases.SaveAuthorUseCase
type SaveAuthorInput = usecases.SaveAuthorInput

var NewListCommentsUseCase = usecases.NewListCommentsUseCase
var NewSaveAuthorUseCase = usecases.NewSaveAuthorUseCase

type Context struct {
	sharedContext *shared.Context
	commentRepo   ports.CommentRepo
}

func NewContext(sharedContext *shared.Context) *Context {
	return &Context{
		sharedContext: sharedContext,
	}
}

func (c *Context) ListCommentsUseCase() *usecases.ListCommentsUseCase {
	return usecases.NewListCommentsUseCase(c.CommentRepo())
}

func (c *Context) SaveAuthorUseCase() *usecases.SaveAuthorUseCase {
	return usecases.NewSaveAuthorUseCase(
		c.CommentRepo(),
		c.sharedContext.TransactionManager(),
		c.sharedContext.IDGenerator(),
	)
}

func (c *Context) CommentRepo() ports.CommentRepo {
	if c.commentRepo == nil {
		c.commentRepo = commentrepo.NewPostgresCommentRepo(c.sharedContext.DB())
	}
	return c.commentRepo
}
