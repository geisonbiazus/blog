package auth

import (
	"context"
	"errors"
	"fmt"
)

type ConfirmOAuth2UseCase struct {
	provider     OAuth2Provider
	stateRepo    StateRepo
	userRepo     UserRepo
	idGen        IDGenerator
	tokenManager TokenManager
}

func NewConfirmOAuth2UseCase(
	provider OAuth2Provider, stateRepo StateRepo, userRepo UserRepo, idGen IDGenerator, tokenManager TokenManager,
) *ConfirmOAuth2UseCase {
	return &ConfirmOAuth2UseCase{
		provider:     provider,
		stateRepo:    stateRepo,
		userRepo:     userRepo,
		idGen:        idGen,
		tokenManager: tokenManager,
	}
}

func (u *ConfirmOAuth2UseCase) Run(ctx context.Context, state, code string) (string, error) {
	providerUser, err := u.processOAuth2Authentication(ctx, state, code)
	if err != nil {
		return "", err
	}

	return u.resolveUserAndGetToken(providerUser)
}

func (u *ConfirmOAuth2UseCase) processOAuth2Authentication(ctx context.Context, state, code string) (ProviderUser, error) {
	err := u.validateAndRemoveState(state)
	if err != nil {
		return ProviderUser{}, err
	}

	return u.getProviderAuthenticatedUser(ctx, code)
}

func (u *ConfirmOAuth2UseCase) validateAndRemoveState(state string) error {
	exists, err := u.stateRepo.Exists(state)
	if err != nil {
		return fmt.Errorf("error checking state on ConfirmOAuth2UseCase: %w", err)
	}

	if !exists {
		return ErrInvalidState
	}

	err = u.stateRepo.Remove(state)
	if err != nil {
		return fmt.Errorf("error authenticating user on ConfirmOAuth2UseCase: %w", err)
	}

	return nil
}

func (u *ConfirmOAuth2UseCase) getProviderAuthenticatedUser(ctx context.Context, code string) (ProviderUser, error) {
	providerUser, err := u.provider.AuthenticatedUser(ctx, code)
	if err != nil {
		return ProviderUser{}, fmt.Errorf("error authenticating user on ConfirmOAuth2UseCase: %w", err)
	}

	return providerUser, nil
}

func (u *ConfirmOAuth2UseCase) resolveUserAndGetToken(providerUser ProviderUser) (string, error) {
	user, err := u.createOrUpdateUser(providerUser)
	if err != nil {
		return "", err
	}

	return u.getAuthenticationToken(user)
}

func (u *ConfirmOAuth2UseCase) createOrUpdateUser(providerUser ProviderUser) (User, error) {
	user, err := u.userRepo.FindUserByProviderUserID(providerUser.ID)

	if errors.Is(err, ErrUserNotFound) {
		return u.createNewUser(providerUser)
	}

	if err != nil {
		return User{}, fmt.Errorf("error finding user on ConfirmOAuth2UseCase: %w", err)
	}

	return u.updateExistingUser(user, providerUser)
}

func (u *ConfirmOAuth2UseCase) createNewUser(providerUser ProviderUser) (User, error) {
	user := User{
		ID:             u.idGen.Generate(),
		ProviderUserID: providerUser.ID,
		Email:          providerUser.Email,
		Name:           providerUser.Name,
		AvatarURL:      providerUser.AvatarURL,
	}

	err := u.userRepo.CreateUser(user)
	if err != nil {
		return User{}, fmt.Errorf("error creatinng user on ConfirmOAuth2UseCase: %w", err)
	}

	return user, nil
}

func (u *ConfirmOAuth2UseCase) updateExistingUser(user User, providerUser ProviderUser) (User, error) {
	user.Email = providerUser.Email
	user.Name = providerUser.Name
	user.AvatarURL = providerUser.AvatarURL

	err := u.userRepo.UpdateUser(user)
	if err != nil {
		return User{}, fmt.Errorf("error updating user on ConfirmOAuth2UseCase: %w", err)
	}

	return user, nil
}

func (u *ConfirmOAuth2UseCase) getAuthenticationToken(user User) (string, error) {
	token, err := u.tokenManager.Encode(user.ID)
	if err != nil {
		return "", fmt.Errorf("error encoding token on ConfirmOAuth2UseCase: %w", err)
	}

	return token, nil
}
