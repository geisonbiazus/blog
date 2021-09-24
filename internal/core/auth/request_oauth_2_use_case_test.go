package auth_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/staterepo/memory"
	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type requestOAuth2UseCaseFixture struct {
	usecase   *auth.RequestOAuth2UseCase
	provider  *OAuth2ProviderSpy
	idGen     *IDGeneratorStub
	stateRepo *memory.InMemoryStateRepo
}

func TestRequestOAuth2UseCase(t *testing.T) {
	setup := func() *requestOAuth2UseCaseFixture {
		stateRepo := memory.NewInMemoryStateRepo()
		idGen := NewIDGeneratorStub()
		provider := NewOAuth2ProviderSpy()
		usecase := auth.NewRequestOAuth2UseCase(provider, idGen, stateRepo)

		return &requestOAuth2UseCaseFixture{
			usecase:   usecase,
			provider:  provider,
			idGen:     idGen,
			stateRepo: stateRepo,
		}
	}

	t.Run("It returns the provider's auth url", func(t *testing.T) {
		f := setup()
		f.provider.ReturnAuthURL = "https://example.com/oauth"

		url, _ := f.usecase.Run()

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

		exists, _ := f.stateRepo.Exists(f.idGen.ReturnID)
		assert.True(t, exists)
	})
}
