package auth

import "fmt"

type RequestOAuth2UseCase struct {
	provider  OAuth2Provider
	idGen     IDGenerator
	stateRepo StateRepo
}

func NewRequestOAuth2UseCase(
	provider OAuth2Provider,
	idGen IDGenerator,
	stateRepo StateRepo,
) *RequestOAuth2UseCase {
	return &RequestOAuth2UseCase{
		provider:  provider,
		idGen:     idGen,
		stateRepo: stateRepo,
	}
}

func (u *RequestOAuth2UseCase) Run() (string, error) {
	state := u.idGen.Generate()

	err := u.stateRepo.AddState(state)
	if err != nil {
		return "", fmt.Errorf("error saving state on RequestOAuth2UseCase: %w", err)
	}

	return u.provider.AuthURL(state), nil
}
