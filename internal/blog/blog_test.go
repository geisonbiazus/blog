package blog_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/app/shared"
	"github.com/geisonbiazus/blog/internal/blog"
	"github.com/geisonbiazus/blog/internal/blog/usecases"
	"github.com/stretchr/testify/suite"
)

type ContextSuite struct {
	suite.Suite
	context *blog.Context
}

func (s *ContextSuite) SetupTest() {
	s.context = blog.NewContext(shared.NewContext())
}
func (s *ContextSuite) TestViewPostUseCase() {
	var usecase *usecases.ViewPostUseCase = s.context.ViewPostUseCase()
	s.NotEqual(nil, usecase)
}

func (s *ContextSuite) TestListPostsUseCase() {
	var usecase *usecases.ListPostsUseCase = s.context.ListPostsUseCase()
	s.NotEqual(nil, usecase)
}

func TestContextSuite(t *testing.T) {
	suite.Run(t, new(ContextSuite))
}
