package auth_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/staterepo/memory"
	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type requestOauth2UseCaseFixture struct {
	usecase   *auth.RequestOauth2UseCase
	provider  *Oauth2ProviderSpy
	idGen     *IDGeneratorStub
	stateRepo *memory.InMemoryStateRepo
}

func TestRequestOauth2UseCase(t *testing.T) {
	setup := func() *requestOauth2UseCaseFixture {
		stateRepo := memory.NewInMemoryStateRepo()
		idGen := NewIDGeneratorStub()
		provider := NewOauth2ProviderSpy()
		usecase := auth.NewRequestOauth2UseCase(provider, idGen, stateRepo)

		return &requestOauth2UseCaseFixture{
			usecase:   usecase,
			provider:  provider,
			idGen:     idGen,
			stateRepo: stateRepo,
		}
	}

	t.Run("It returns the provider's auth url", func(t *testing.T) {
		f := setup()
		f.provider.ReturnAuthURL = "https://example.com/oauth"

		url := f.usecase.Run()

		assert.Equal(t, url, f.provider.ReturnAuthURL)
	})

	t.Run("It generates a random state and sends it to the provider", func(t *testing.T) {
		f := setup()
		f.idGen.ReturnID = "random_state"

		f.usecase.Run()

		assert.Equal(t, f.idGen.ReturnID, f.provider.ReceivedState)
	})

	t.Run("It persists the generated state", func(t *testing.T) {
		f := setup()
		f.idGen.ReturnID = "random_state"

		f.usecase.Run()

		assert.True(t, f.stateRepo.StateExists(f.idGen.ReturnID))
	})
}
