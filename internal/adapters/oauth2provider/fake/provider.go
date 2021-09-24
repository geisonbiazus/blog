package fake

import (
	"context"
	"fmt"

	"github.com/geisonbiazus/blog/internal/adapters/idgenerator/uuid"
	"github.com/geisonbiazus/blog/internal/core/auth"
)

type Provider struct {
	BaseURL string
	State   string
	Code    string
}

func NewProvider(baseURL string) *Provider {
	return &Provider{
		BaseURL: baseURL,
	}
}

func (p *Provider) AuthURL(state string) string {
	p.State = state
	p.Code = uuid.NewGenerator().Generate()
	return fmt.Sprintf("%s/login/github/confirm?state=%s&code=%s", p.BaseURL, p.State, p.Code)
}

func (p *Provider) AuthenticatedUser(ctx context.Context, code string) (auth.ProviderUser, error) {
	return auth.ProviderUser{
		ID:        "userID",
		Email:     "user@example.com",
		Name:      "User",
		AvatarURL: "http://example.com/avatar",
	}, nil
}
