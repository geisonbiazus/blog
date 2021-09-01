package auth_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestRequestOauth2UseCase(t *testing.T) {
	t.Run("it requests and returns the provider's auth url", func(t *testing.T) {
		provider := NewOauth2ProviderSpy()
		usecase := auth.NewRequestOauth2UseCase(provider)

		provider.ReturnAuthURL = "https://example.com/oauth"

		url := usecase.Run()
		assert.Equal(t, url, provider.ReturnAuthURL)
	})
}
