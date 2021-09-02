package auth

type RequestOauth2UseCase struct {
	provider  Oauth2Provider
	idGen     IDGenerator
	stateRepo StateRepo
}

func NewRequestOauth2UseCase(provider Oauth2Provider, idGen IDGenerator, stateRepo StateRepo) *RequestOauth2UseCase {
	return &RequestOauth2UseCase{provider: provider, idGen: idGen, stateRepo: stateRepo}
}

func (u *RequestOauth2UseCase) Run() string {
	state := u.idGen.Generate()
	u.stateRepo.AddState(state)

	return u.provider.AuthURL(state)
}
