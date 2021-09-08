package auth_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type confirmOauth2UseCaseFixture struct {
	usecase *auth.ConfimOauth2UseCase
}

func TestConfirmOauth2UseCase(t *testing.T) {
	setup := func() *confirmOauth2UseCaseFixture {
		usecase := auth.NewConfimOauth2UseCase()
		return &confirmOauth2UseCaseFixture{usecase: usecase}
	}

	t.Run("It returns error when state is not found", func(t *testing.T) {
		f := setup()

		err := f.usecase.Run("state", "code")

		assert.Equal(t, err, auth.ErrInvalidState)
	})
}
