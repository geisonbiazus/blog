package oauth2provider

import (
	"github.com/geisonbiazus/blog/internal/adapters/oauth2provider/fake"
	"github.com/geisonbiazus/blog/internal/adapters/oauth2provider/github"
)

func NewFakeProvider(baseURL string) *fake.Provider {
	return fake.NewProvider(baseURL)
}

func NewGithubProvider(clientID, clientSecret string) *github.Provider {
	return github.NewProvider(clientID, clientSecret)
}
