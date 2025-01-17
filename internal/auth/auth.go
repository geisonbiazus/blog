package auth

import (
	"github.com/geisonbiazus/blog/internal/app/shared"
	"github.com/geisonbiazus/blog/internal/auth/adapters/oauth2provider"
	"github.com/geisonbiazus/blog/internal/auth/adapters/staterepo"
	"github.com/geisonbiazus/blog/internal/auth/adapters/tokenencoder"
	"github.com/geisonbiazus/blog/internal/auth/adapters/userrepo"
	"github.com/geisonbiazus/blog/internal/auth/entities"
	"github.com/geisonbiazus/blog/internal/auth/events"
	"github.com/geisonbiazus/blog/internal/auth/ports"
	"github.com/geisonbiazus/blog/internal/auth/usecases"
	"github.com/geisonbiazus/blog/pkg/env"
)

type ProviderUser = entities.ProviderUser
type User = entities.User

type OAuth2Provider = ports.OAuth2Provider
type StateRepo = ports.StateRepo
type UserRepo = ports.UserRepo
type TokenEncoder = ports.TokenEncoder

const UserCreatedEvent = events.UserCreatedEvent
const UserUpdatedEvent = events.UserUpdatedEvent

var ErrInvalidState = entities.ErrInvalidState
var ErrUserNotFound = entities.ErrUserNotFound
var ErrTokenExpired = entities.ErrTokenExpired

type Context struct {
	GitHubClientID     string
	GitHubClientSecret string
	AuthTokenSecret    string

	sharedContext *shared.Context

	stateRepo ports.StateRepo
	userRepo  ports.UserRepo
}

func NewContext(sharedContext *shared.Context) *Context {
	return &Context{
		sharedContext: sharedContext,

		GitHubClientID:     env.GetString("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: env.GetString("GITHUB_CLIENT_SECRET", ""),
		AuthTokenSecret:    env.GetString("AUTH_TOKEN_SECRET", ""),
	}
}

// Use cases

func (c *Context) RequestOAuth2UseCase() *usecases.RequestOAuth2UseCase {
	return usecases.NewRequestOAuth2UseCase(c.OAuth2Provider(), c.sharedContext.IDGenerator(), c.StateRepo())
}

func (c *Context) ConfirmOAuth2UseCase() *usecases.ConfirmOAuth2UseCase {
	return usecases.NewConfirmOAuth2UseCase(
		c.OAuth2Provider(),
		c.StateRepo(),
		c.UserRepo(),
		c.sharedContext.IDGenerator(),
		c.TokenEncoder(),
		c.sharedContext.TransactionManager(),
		c.sharedContext.Publisher())
}

// Adapters

func (c *Context) OAuth2Provider() ports.OAuth2Provider {
	if c.sharedContext.IsTest() {
		return c.FakeOAuth2Provider()
	}
	return c.GithubOAuth2Provider()
}

func (c *Context) GithubOAuth2Provider() ports.OAuth2Provider {
	return oauth2provider.NewGithubProvider(c.GitHubClientID, c.GitHubClientSecret)
}

func (c *Context) FakeOAuth2Provider() ports.OAuth2Provider {
	return oauth2provider.NewFakeProvider(c.sharedContext.BaseURL)
}

func (c *Context) StateRepo() ports.StateRepo {
	if c.stateRepo == nil {
		c.stateRepo = staterepo.NewMemoryStateRepo()
	}
	return c.stateRepo
}

func (c *Context) UserRepo() ports.UserRepo {
	if c.userRepo == nil {
		c.userRepo = userrepo.NewPostgresUserRepo(c.sharedContext.DB())
	}
	return c.userRepo
}

func (c *Context) TokenEncoder() ports.TokenEncoder {
	return tokenencoder.NewJWTTokenEncoder(c.AuthTokenSecret)
}
