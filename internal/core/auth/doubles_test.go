package auth_test

import (
	"context"
	"time"

	"github.com/geisonbiazus/blog/internal/core/auth"
)

type OAuth2ProviderSpy struct {
	ReturnAuthURL string
	ReceivedState string

	AuthenticatedUserReturnProviderUser auth.ProviderUser
	AuthenticatedUserReturnError        error
	AuthenticatedUserReceivedContext    context.Context
	AuthenticatedUserReceivedCode       string
}

func NewOAuth2ProviderSpy() *OAuth2ProviderSpy {
	return &OAuth2ProviderSpy{ReturnAuthURL: "https://example.com/oauth"}
}

func (p *OAuth2ProviderSpy) AuthURL(state string) string {
	p.ReceivedState = state
	return p.ReturnAuthURL
}

func (p *OAuth2ProviderSpy) AuthenticatedUser(ctx context.Context, code string) (auth.ProviderUser, error) {
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

type TokenEncoderSpy struct {
	EncodeReturnToken       string
	EncodeReturnError       error
	EncodeReceivedValue     string
	EncodeReceivedExpiresIn time.Duration
}

func NewTokenEncoderSpy() *TokenEncoderSpy {
	return &TokenEncoderSpy{}
}

func (m *TokenEncoderSpy) Encode(value string, expiresIn time.Duration) (string, error) {
	m.EncodeReceivedValue = value
	m.EncodeReceivedExpiresIn = expiresIn
	return m.EncodeReturnToken, m.EncodeReturnError
}
