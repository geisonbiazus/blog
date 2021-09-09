package github

import (
	"context"

	"github.com/geisonbiazus/blog/internal/core/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Provider struct {
	config oauth2.Config
}

func NewProvider(clientID, clientSecret string) *Provider {
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     github.Endpoint,
	}

	return &Provider{config: config}
}

func (p *Provider) AuthURL(state string) string {
	return p.config.AuthCodeURL(state)
}

func (p *Provider) AuthenticatedUser(ctx context.Context, code string) (auth.ProviderUser, error) {
	return auth.ProviderUser{}, nil
}
