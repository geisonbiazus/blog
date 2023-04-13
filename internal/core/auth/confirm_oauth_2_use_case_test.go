package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	fakepublisher "github.com/geisonbiazus/blog/internal/adapters/publisher/fake"
	staterepo "github.com/geisonbiazus/blog/internal/adapters/staterepo/memory"
	"github.com/geisonbiazus/blog/internal/adapters/transactionmanager/fake"
	userrepo "github.com/geisonbiazus/blog/internal/adapters/userrepo/memory"
	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/internal/core/shared"
	"github.com/stretchr/testify/assert"
)

type confirmOAuth2UseCaseFixture struct {
	usecase      *auth.ConfirmOAuth2UseCase
	provider     *OAuth2ProviderSpy
	stateRepo    *staterepo.StateRepo
	userRepo     *userrepo.UserRepo
	idGen        *IDGeneratorStub
	tokenEncoder *TokenEncoderSpy
	publisher    *fakepublisher.Publisher
	ctx          context.Context
}

func TestConfirmOAuth2UseCase(t *testing.T) {
	code := "code"
	state := "state"

	providerUser := auth.ProviderUser{
		ID:        "provider_id",
		Email:     "user@example.com",
		Name:      "name",
		AvatarURL: "http://example.com/avatar.png",
	}

	setup := func() *confirmOAuth2UseCaseFixture {
		provider := NewOAuth2ProviderSpy()
		stateRepo := staterepo.NewStateRepo()
		userRepo := userrepo.NewUserRepo()
		idGen := NewIDGeneratorStub()
		tokenEncoder := NewTokenEncoderSpy()
		txManager := fake.NewTransactionManager()
		publisher := fakepublisher.NewPublisher()
		usecase := auth.NewConfirmOAuth2UseCase(provider, stateRepo, userRepo, idGen, tokenEncoder, txManager, publisher)
		return &confirmOAuth2UseCaseFixture{
			usecase:      usecase,
			provider:     provider,
			stateRepo:    stateRepo,
			userRepo:     userRepo,
			idGen:        idGen,
			tokenEncoder: tokenEncoder,
			publisher:    publisher,
			ctx:          context.Background(),
		}
	}

	t.Run("It returns error when state is not found", func(t *testing.T) {
		f := setup()

		_, err := f.usecase.Run(f.ctx, state, code)

		assert.Equal(t, err, auth.ErrInvalidState)
	})

	t.Run("It removes state when it is found", func(t *testing.T) {
		f := setup()

		f.stateRepo.AddState(state)

		f.usecase.Run(f.ctx, state, code)

		exists, _ := f.stateRepo.Exists(state)
		assert.False(t, exists)
	})

	t.Run("It returns error when it fails to get the authenticated user from the provider", func(t *testing.T) {
		f := setup()

		f.stateRepo.AddState(state)

		authError := errors.New("authentication error")
		f.provider.AuthenticatedUserReturnError = authError

		_, err := f.usecase.Run(f.ctx, state, code)

		assert.Equal(t, f.ctx, f.provider.AuthenticatedUserReceivedContext)
		assert.Equal(t, code, f.provider.AuthenticatedUserReceivedCode)
		assert.Error(t, authError, err)
	})

	t.Run("It creates a user when authentication is successful", func(t *testing.T) {
		f := setup()

		f.stateRepo.AddState(state)
		f.idGen.ReturnID = "generatedID"
		f.provider.AuthenticatedUserReturnProviderUser = providerUser

		f.usecase.Run(f.ctx, state, code)

		user := auth.User{
			ID:             f.idGen.ReturnID,
			ProviderUserID: providerUser.ID,
			Email:          providerUser.Email,
			Name:           providerUser.Name,
			AvatarURL:      providerUser.AvatarURL,
		}

		createdUser, _ := f.userRepo.FindUserByEmail(f.ctx, providerUser.Email)

		assert.Equal(t, user, createdUser)

		assert.Equal(t, shared.Event{
			Type:       auth.UserCreatedEvent,
			OccurredOn: f.publisher.LastEvent().OccurredOn,
			Payload: map[string]interface{}{
				"ID":        user.ID,
				"Email":     user.Email,
				"Name":      user.Name,
				"AvatarURL": user.AvatarURL,
			},
		}, f.publisher.LastEvent())
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

		f.userRepo.CreateUser(f.ctx, user)

		f.usecase.Run(f.ctx, state, code)

		expctedUser := auth.User{
			ID:             user.ID,
			ProviderUserID: providerUser.ID,
			Email:          providerUser.Email,
			Name:           providerUser.Name,
			AvatarURL:      providerUser.AvatarURL,
		}

		createdUser, _ := f.userRepo.FindUserByEmail(f.ctx, providerUser.Email)

		assert.Equal(t, expctedUser, createdUser)

		assert.Equal(t, shared.Event{
			Type:       auth.UserUpdatedEvent,
			OccurredOn: f.publisher.LastEvent().OccurredOn,
			Payload: map[string]interface{}{
				"ID":        expctedUser.ID,
				"Email":     expctedUser.Email,
				"Name":      expctedUser.Name,
				"AvatarURL": expctedUser.AvatarURL,
			},
		}, f.publisher.LastEvent())
	})

	t.Run("It authenticates the user and returns the authentication token", func(t *testing.T) {
		f := setup()

		f.stateRepo.AddState(state)
		f.idGen.ReturnID = "generatedID"
		f.provider.AuthenticatedUserReturnProviderUser = providerUser
		f.tokenEncoder.EncodeReturnToken = "expectedToken"

		token, _ := f.usecase.Run(f.ctx, state, code)

		assert.Equal(t, f.idGen.ReturnID, f.tokenEncoder.EncodeReceivedValue)
		assert.Equal(t, 24*time.Hour, f.tokenEncoder.EncodeReceivedExpiresIn)
		assert.Equal(t, f.tokenEncoder.EncodeReturnToken, token)
	})

	t.Run("It returns error it fails to generate token", func(t *testing.T) {
		f := setup()

		f.stateRepo.AddState(state)
		f.idGen.ReturnID = "generatedID"
		f.provider.AuthenticatedUserReturnProviderUser = providerUser
		f.tokenEncoder.EncodeReturnError = errors.New("error encoding")

		_, err := f.usecase.Run(f.ctx, state, code)

		assert.Error(t, f.tokenEncoder.EncodeReturnError, err)
	})
}
