package auth

import "fmt"

type RequestOauth2UseCase struct {
	provider  Oauth2Provider
	idGen     IDGenerator
	stateRepo StateRepo
}

func NewRequestOauth2UseCase(provider Oauth2Provider, idGen IDGenerator, stateRepo StateRepo) *RequestOauth2UseCase {
	return &RequestOauth2UseCase{provider: provider, idGen: idGen, stateRepo: stateRepo}
}

func (u *RequestOauth2UseCase) Run() (string, error) {
	state := u.idGen.Generate()

	err := u.stateRepo.AddState(state)
	if err != nil {
		return "", fmt.Errorf("error saving state on RequestOauth2UseCase: %w", err)
	}

	return u.provider.AuthURL(state), nil
}
