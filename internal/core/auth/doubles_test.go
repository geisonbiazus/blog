package auth_test

import (
	"context"

	"github.com/geisonbiazus/blog/internal/core/auth"
)

type Oauth2ProviderSpy struct {
	ReturnAuthURL string
	ReceivedState string

	AuthenticatedUserReturnProviderUser auth.ProviderUser
	AuthenticatedUserReturnError        error
	AuthenticatedUserReceivedContext    context.Context
	AuthenticatedUserReceivedCode       string
}

func NewOauth2ProviderSpy() *Oauth2ProviderSpy {
	return &Oauth2ProviderSpy{ReturnAuthURL: "https://example.com/oauth"}
}

func (p *Oauth2ProviderSpy) AuthURL(state string) string {
	p.ReceivedState = state
	return p.ReturnAuthURL
}

func (p *Oauth2ProviderSpy) AuthenticatedUser(ctx context.Context, code string) (auth.ProviderUser, error) {
	p.AuthenticatedUserReceivedContext = ctx
	p.AuthenticatedUserReceivedCode = code
	return p.AuthenticatedUserReturnProviderUser, p.AuthenticatedUserReturnError
}

type IDGeneratorStub struct {
	ReturnID string
}

func NewIDGeneratorStub() *IDGeneratorStub {
	return &IDGeneratorStub{}
}

func (g *IDGeneratorStub) Generate() string {
	return g.ReturnID
}

type TokenManagerSpy struct {
	EncodeReturnToken    string
	EncodeReturnError    error
	EncodeReceivedUserID string
}

func NewTokenManagerSpy() *TokenManagerSpy {
	return &TokenManagerSpy{}
}

func (m *TokenManagerSpy) Encode(userID string) (string, error) {
	m.EncodeReceivedUserID = userID
	return m.EncodeReturnToken, m.EncodeReturnError
}
