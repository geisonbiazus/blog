package github_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/oauth2provider/github"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	t.Run("AuthURL", func(t *testing.T) {
		t.Run("It generates github auth URL", func(t *testing.T) {
			clientID := "client-id"
			clientSecret := "client-secret"

			provider := github.NewProvider(clientID, clientSecret)
			assert.Equal(t, "https://github.com/login/oauth/authorize?client_id=client-id&response_type=code&state=state", provider.AuthURL("state"))
		})
	})
}
