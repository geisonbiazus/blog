package blog_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/blog"
	"github.com/geisonbiazus/blog/internal/blog/usecases"
	"github.com/geisonbiazus/blog/pkg/caching"
	"github.com/stretchr/testify/suite"
)

type ContextSuite struct {
	suite.Suite
	context *blog.Context
}

func (s *ContextSuite) SetupTest() {
	cacheFn := func() caching.Cache { return caching.NewNullCache() }
	s.context = blog.NewContext(cacheFn)
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
