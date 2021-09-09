package auth

import (
	"context"
	"errors"
	"fmt"
)

type ConfirmOauth2UseCase struct {
	provider     Oauth2Provider
	stateRepo    StateRepo
	userRepo     UserRepo
	idGen        IDGenerator
	tokenManager TokenManager
}

func NewConfirmOauth2UseCase(
	provider Oauth2Provider, stateRepo StateRepo, userRepo UserRepo, idGen IDGenerator, tokenManager TokenManager,
) *ConfirmOauth2UseCase {
	return &ConfirmOauth2UseCase{
		provider:     provider,
		stateRepo:    stateRepo,
		userRepo:     userRepo,
		idGen:        idGen,
		tokenManager: tokenManager,
	}
}

func (u *ConfirmOauth2UseCase) Run(ctx context.Context, state, code string) (string, error) {
	err := u.validateAndRemoveState(state)
	if err != nil {
		return "", err
	}

	providerUser, err := u.getProviderAuthenticatedUser(ctx, code)
	if err != nil {
		return "", err
	}

	user, err := u.createOrUpdateUser(providerUser)
	if err != nil {
		return "", err
	}

	token, err := u.authenticateUser(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *ConfirmOauth2UseCase) validateAndRemoveState(state string) error {
	exists, err := u.stateRepo.Exists(state)
	if err != nil {
		return fmt.Errorf("error checking state on ConfirmOauth2UseCase: %w", err)
	}

	if !exists {
		return ErrInvalidState
	}

	err = u.stateRepo.Remove(state)
	if err != nil {
		return fmt.Errorf("error authenticating user on ConfirmOauth2UseCase: %w", err)
	}

	return nil
}

func (u *ConfirmOauth2UseCase) getProviderAuthenticatedUser(ctx context.Context, code string) (ProviderUser, error) {
	providerUser, err := u.provider.AuthenticatedUser(ctx, code)
	if err != nil {
		return ProviderUser{}, fmt.Errorf("error authenticating user on ConfirmOauth2UseCase: %w", err)
	}

	return providerUser, nil
}

func (u *ConfirmOauth2UseCase) createOrUpdateUser(providerUser ProviderUser) (User, error) {
	user, err := u.userRepo.FindUserByProviderUserID(providerUser.ID)

	if errors.Is(err, ErrUserNotFound) {
		return u.createNewUser(providerUser)
	}

	if err != nil {
		return User{}, fmt.Errorf("error finding user on ConfirmOauth2UseCase: %w", err)
	}

	return u.updateExistingUser(user, providerUser)
}

func (u *ConfirmOauth2UseCase) createNewUser(providerUser ProviderUser) (User, error) {
	user := User{
		ID:             u.idGen.Generate(),
		ProviderUserID: providerUser.ID,
		Email:          providerUser.Email,
		Name:           providerUser.Name,
		AvatarURL:      providerUser.AvatarURL,
	}

	err := u.userRepo.CreateUser(user)
	if err != nil {
		return User{}, fmt.Errorf("error creatinng user on ConfirmOauth2UseCase: %w", err)
	}
	return user, nil
}

func (u *ConfirmOauth2UseCase) updateExistingUser(user User, providerUser ProviderUser) (User, error) {
	user.Email = providerUser.Email
	user.Name = providerUser.Name
	user.AvatarURL = providerUser.AvatarURL

	err := u.userRepo.UpdateUser(user)
	if err != nil {
		return User{}, fmt.Errorf("error updating user on ConfirmOauth2UseCase: %w", err)
	}

	return user, nil
}

func (u *ConfirmOauth2UseCase) authenticateUser(user User) (string, error) {
	token, err := u.tokenManager.Encode(user.ID)
	if err != nil {
		return "", fmt.Errorf("error encoding token on ConfirmOauth2UseCase: %w", err)
	}

	return token, nil
}
