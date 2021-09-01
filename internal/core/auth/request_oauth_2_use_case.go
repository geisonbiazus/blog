package auth

type RequestOauth2UseCase struct {
	provider Oauth2Provider
}

func NewRequestOauth2UseCase(provider Oauth2Provider) *RequestOauth2UseCase {
	return &RequestOauth2UseCase{provider: provider}
}

func (u *RequestOauth2UseCase) Run() string {
	return u.provider.AuthURL()
}
