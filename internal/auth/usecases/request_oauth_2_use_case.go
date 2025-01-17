package usecases

import (
	"fmt"

	"github.com/geisonbiazus/blog/internal/auth/ports"
	"github.com/geisonbiazus/blog/pkg/gen"
)

type RequestOAuth2UseCase struct {
	provider  ports.OAuth2Provider
	idGen     gen.Generator
	stateRepo ports.StateRepo
}

func NewRequestOAuth2UseCase(
	provider ports.OAuth2Provider,
	idGen gen.Generator,
	stateRepo ports.StateRepo,
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
