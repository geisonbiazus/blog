package discussion_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/app/shared"
	"github.com/geisonbiazus/blog/internal/discussion"
	"github.com/geisonbiazus/blog/internal/discussion/usecases"
	"github.com/stretchr/testify/suite"
)

type ContextSuite struct {
	suite.Suite
	context *discussion.Context
}

func (s *ContextSuite) SetupTest() {
	s.context = discussion.NewContext(shared.NewContext())
}

func (s *ContextSuite) TestListCommentsUseCase() {
	var usecase *usecases.ListCommentsUseCase = s.context.ListCommentsUseCase()
	s.NotNil(usecase)
}

func (s *ContextSuite) TestSaveAuthorUseCase() {
	var usecase *usecases.SaveAuthorUseCase = s.context.SaveAuthorUseCase()
	s.NotNil(usecase)
}

func TestContextSuite(t *testing.T) {
	suite.Run(t, new(ContextSuite))
}
