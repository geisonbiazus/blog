package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/geisonbiazus/blog/internal/auth/entities"
	"github.com/geisonbiazus/blog/internal/auth/events"
	"github.com/geisonbiazus/blog/internal/auth/ports"
	"github.com/geisonbiazus/blog/pkg/eventing"
	"github.com/geisonbiazus/blog/pkg/gen"
	"github.com/geisonbiazus/blog/pkg/transaction"
)

type ConfirmOAuth2UseCase struct {
	provider     ports.OAuth2Provider
	stateRepo    ports.StateRepo
	userRepo     ports.UserRepo
	idGen        gen.Generator
	tokenEncoder ports.TokenEncoder
	txManager    transaction.Manager
	publisher    eventing.Publisher
}

func NewConfirmOAuth2UseCase(
	provider ports.OAuth2Provider,
	stateRepo ports.StateRepo,
	userRepo ports.UserRepo,
	idGen gen.Generator,
	tokenEncoder ports.TokenEncoder,
	txManager transaction.Manager,
	publisher eventing.Publisher,
) *ConfirmOAuth2UseCase {
	return &ConfirmOAuth2UseCase{
		provider:     provider,
		stateRepo:    stateRepo,
		userRepo:     userRepo,
		idGen:        idGen,
		tokenEncoder: tokenEncoder,
		txManager:    txManager,
		publisher:    publisher,
	}
}

func (u *ConfirmOAuth2UseCase) Run(ctx context.Context, state, code string) (token string, err error) {
	err = u.txManager.Transaction(ctx, func(ctx context.Context) error {
		token, err = u.run(ctx, state, code)
		return err
	})
	return
}

func (u *ConfirmOAuth2UseCase) run(ctx context.Context, state, code string) (string, error) {
	providerUser, err := u.processOAuth2Authentication(ctx, state, code)
	if err != nil {
		return "", err
	}

	return u.resolveUserAndGetToken(ctx, providerUser)
}

func (u *ConfirmOAuth2UseCase) processOAuth2Authentication(ctx context.Context, state, code string) (entities.ProviderUser, error) {
	err := u.validateAndRemoveState(state)
	if err != nil {
		return entities.ProviderUser{}, err
	}

	return u.getProviderAuthenticatedUser(ctx, code)
}

func (u *ConfirmOAuth2UseCase) validateAndRemoveState(state string) error {
	exists, err := u.stateRepo.Exists(state)
	if err != nil {
		return fmt.Errorf("error checking state on ConfirmOAuth2UseCase: %w", err)
	}

	if !exists {
		return entities.ErrInvalidState
	}

	err = u.stateRepo.Remove(state)
	if err != nil {
		return fmt.Errorf("error authenticating user on ConfirmOAuth2UseCase: %w", err)
	}

	return nil
}

func (u *ConfirmOAuth2UseCase) getProviderAuthenticatedUser(ctx context.Context, code string) (entities.ProviderUser, error) {
	providerUser, err := u.provider.AuthenticatedUser(ctx, code)
	if err != nil {
		return entities.ProviderUser{}, fmt.Errorf("error authenticating user on ConfirmOAuth2UseCase: %w", err)
	}

	return providerUser, nil
}

func (u *ConfirmOAuth2UseCase) resolveUserAndGetToken(ctx context.Context, providerUser entities.ProviderUser) (string, error) {
	user, err := u.createOrUpdateUser(ctx, providerUser)
	if err != nil {
		return "", err
	}

	return u.getAuthenticationToken(user)
}

func (u *ConfirmOAuth2UseCase) createOrUpdateUser(ctx context.Context, providerUser entities.ProviderUser) (entities.User, error) {
	user, err := u.userRepo.FindUserByProviderUserID(ctx, providerUser.ID)

	if errors.Is(err, entities.ErrUserNotFound) {
		return u.createNewUser(ctx, providerUser)
	}

	if err != nil {
		return entities.User{}, fmt.Errorf("error finding user on ConfirmOAuth2UseCase: %w", err)
	}

	return u.updateExistingUser(ctx, user, providerUser)
}

func (u *ConfirmOAuth2UseCase) createNewUser(ctx context.Context, providerUser entities.ProviderUser) (entities.User, error) {
	user := entities.User{
		ID:             u.idGen.Generate(),
		ProviderUserID: providerUser.ID,
		Email:          providerUser.Email,
		Name:           providerUser.Name,
		AvatarURL:      providerUser.AvatarURL,
	}

	err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return entities.User{}, fmt.Errorf("error creatinng user on ConfirmOAuth2UseCase: %w", err)
	}

	err = u.publisher.Publish(events.NewUserCreatedEvent(user))
	if err != nil {
		return entities.User{}, fmt.Errorf("error publishing event on ConfirmOAuth2UseCase: %w", err)
	}

	return user, nil
}

func (u *ConfirmOAuth2UseCase) updateExistingUser(ctx context.Context, user entities.User, providerUser entities.ProviderUser) (entities.User, error) {
	user.Email = providerUser.Email
	user.Name = providerUser.Name
	user.AvatarURL = providerUser.AvatarURL

	err := u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return entities.User{}, fmt.Errorf("error updating user on ConfirmOAuth2UseCase: %w", err)
	}

	err = u.publisher.Publish(events.NewUserUpdatedEvent(user))
	if err != nil {
		return entities.User{}, fmt.Errorf("error publishing event on ConfirmOAuth2UseCase: %w", err)
	}

	return user, nil
}

const TokenExpiration = 24 * time.Hour

func (u *ConfirmOAuth2UseCase) getAuthenticationToken(user entities.User) (string, error) {
	token, err := u.tokenEncoder.Encode(user.ID, TokenExpiration)
	if err != nil {
		return "", fmt.Errorf("error encoding token on ConfirmOAuth2UseCase: %w", err)
	}

	return token, nil
}
