package blog_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/app/shared"
	"github.com/geisonbiazus/blog/internal/blog"
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
	var usecase blog.ViewPostUseCase = s.context.ViewPostUseCase()
	s.NotNil(usecase)
}

func (s *ContextSuite) TestListPostsUseCase() {
	var usecase blog.ListPostsUseCase = s.context.ListPostsUseCase()
	s.NotNil(usecase)
}

func TestContextSuite(t *testing.T) {
	suite.Run(t, new(ContextSuite))
}
