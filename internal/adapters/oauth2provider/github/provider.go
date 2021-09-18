package github

import (
	"context"
	"fmt"
	"net/http"

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
	httpClient, err := p.exhangeTokenAndGetClient(ctx, code)
	if err != nil {
		return auth.ProviderUser{}, err
	}

	return NewClient(httpClient).GetAuthenticatedUser()
}

func (p *Provider) exhangeTokenAndGetClient(ctx context.Context, code string) (*http.Client, error) {
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("error exchanging token on github.Provider: %w", err)
	}

	tokenSource := p.config.TokenSource(ctx, token)
	return oauth2.NewClient(ctx, tokenSource), nil
}

type githubResponse struct {
	Data struct {
		Viewer struct {
			ID        string `json:"id"`
			Email     string `json:"email"`
			Name      string `json:"name"`
			AvatarURL string `json:"avatarUrl"`
		} `json:"viewer"`
	} `json:"data"`
}
