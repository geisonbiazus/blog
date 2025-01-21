package auth_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/app/shared"
	"github.com/geisonbiazus/blog/internal/auth"
	"github.com/stretchr/testify/suite"
)

type ContextSuite struct {
	suite.Suite
	context *auth.Context
}

func (s *ContextSuite) SetupTest() {
	sharedContext := shared.NewContext()
	s.context = auth.NewContext(sharedContext)
}

func (s *ContextSuite) TestRequestOAuth2UseCase() {
	var usecase auth.RequestOAuth2UseCase = s.context.RequestOAuth2UseCase()
	s.NotNil(usecase)
}

func (s *ContextSuite) TestConfirmOAuth2UseCase() {
	var usecase auth.ConfirmOAuth2UseCase = s.context.ConfirmOAuth2UseCase()
	s.NotNil(usecase)
}
func TestContextSuite(t *testing.T) {
	suite.Run(t, new(ContextSuite))
}
