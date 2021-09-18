package auth_test

import (
	"context"
	"errors"
	"testing"

	staterepo "github.com/geisonbiazus/blog/internal/adapters/staterepo/memory"
	userrepo "github.com/geisonbiazus/blog/internal/adapters/userrepo/memory"
	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/pkg/assert"
)

type confirmOauth2UseCaseFixture struct {
	usecase      *auth.ConfirmOauth2UseCase
	provider     *Oauth2ProviderSpy
	stateRepo    *staterepo.InMemoryStateRepo
	userRepo     *userrepo.UserRepo
	idGen        *IDGeneratorStub
	tokenManager *TokenManagerSpy
	context      context.Context
}

func TestConfirmOauth2UseCase(t *testing.T) {
	code := "code"
	state := "state"

	providerUser := auth.ProviderUser{
		ID:        "provider_id",
		Email:     "user@example.com",
		Name:      "name",
		AvatarURL: "http://example.com/avatar.png",
	}

	setup := func() *confirmOauth2UseCaseFixture {
		provider := NewOauth2ProviderSpy()
		stateRepo := staterepo.NewInMemoryStateRepo()
		userRepo := userrepo.NewUserRepo()
		idGen := NewIDGeneratorStub()
		tokenManager := NewTokenManagerSpy()
		usecase := auth.NewConfirmOauth2UseCase(provider, stateRepo, userRepo, idGen, tokenManager)
		return &confirmOauth2UseCaseFixture{
			usecase:      usecase,
			provider:     provider,
			stateRepo:    stateRepo,
			userRepo:     userRepo,
			idGen:        idGen,
			tokenManager: tokenManager,
			context:      context.Background(),
		}
	}

	t.Run("It returns error when state is not found", func(t *testing.T) {
		f := setup()

		_, err := f.usecase.Run(f.context, state, code)

		assert.Equal(t, err, auth.ErrInvalidState)
	})

	t.Run("It removes state when it is found", func(t *testing.T) {
		f := setup()

		f.stateRepo.AddState(state)

		f.usecase.Run(f.context, state, code)

		exists, _ := f.stateRepo.Exists(state)
		assert.False(t, exists)
	})

	t.Run("It returns error when it fails to get the authenticated user from the provider", func(t *testing.T) {
		f := setup()

		f.stateRepo.AddState(state)

		authError := errors.New("authentication error")
		f.provider.AuthenticatedUserReturnError = authError

		_, err := f.usecase.Run(f.context, state, code)

		assert.Equal(t, f.context, f.provider.AuthenticatedUserReceivedContext)
		assert.Equal(t, code, f.provider.AuthenticatedUserReceivedCode)
		assert.Error(t, authError, err)
	})

	t.Run("It creates a user when authentication is successful", func(t *testing.T) {
		f := setup()

		f.stateRepo.AddState(state)
		f.idGen.ReturnID = "generatedID"
		f.provider.AuthenticatedUserReturnProviderUser = providerUser

		f.usecase.Run(f.context, state, code)

		user := auth.User{
			ID:             f.idGen.ReturnID,
			ProviderUserID: providerUser.ID,
			Email:          providerUser.Email,
			Name:           providerUser.Name,
			AvatarURL:      providerUser.AvatarURL,
		}

		createdUser, _ := f.userRepo.FindUserByEmail(providerUser.Email)

		assert.Equal(t, user, createdUser)
	})

	t.Run("It updates user data when authentication is successful and the user already exists", func(t *testing.T) {
		f := setup()

		f.stateRepo.AddState(state)
		f.idGen.ReturnID = "generatedID"
		f.provider.AuthenticatedUserReturnProviderUser = providerUser

		user := auth.User{
			ID:             "existing user ID",
			ProviderUserID: providerUser.ID,
			Email:          "previous.email@example.com",
			Name:           "previous name",
			AvatarURL:      "http://example.com/previous_avatar.png",
		}

		f.userRepo.CreateUser(user)

		f.usecase.Run(f.context, state, code)

		expctedUser := auth.User{
			ID:             user.ID,
			ProviderUserID: providerUser.ID,
			Email:          providerUser.Email,
			Name:           providerUser.Name,
			AvatarURL:      providerUser.AvatarURL,
		}

		createdUser, _ := f.userRepo.FindUserByEmail(providerUser.Email)

		assert.Equal(t, expctedUser, createdUser)
	})

	t.Run("It authenticates the user and returns the authentication token", func(t *testing.T) {
		f := setup()

		f.stateRepo.AddState(state)
		f.idGen.ReturnID = "generatedID"
		f.provider.AuthenticatedUserReturnProviderUser = providerUser
		f.tokenManager.EncodeReturnToken = "expectedToken"

		token, _ := f.usecase.Run(f.context, state, code)

		assert.Equal(t, f.idGen.ReturnID, f.tokenManager.EncodeReceivedUserID)
		assert.Equal(t, f.tokenManager.EncodeReturnToken, token)
	})

	t.Run("It returns error it fails to generate token", func(t *testing.T) {
		f := setup()

		f.stateRepo.AddState(state)
		f.idGen.ReturnID = "generatedID"
		f.provider.AuthenticatedUserReturnProviderUser = providerUser
		f.tokenManager.EncodeReturnError = errors.New("error encoding")

		_, err := f.usecase.Run(f.context, state, code)

		assert.Error(t, f.tokenManager.EncodeReturnError, err)
	})
}