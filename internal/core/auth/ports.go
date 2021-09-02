package auth

type Oauth2Provider interface {
	AuthURL(state string) string
}

type IDGenerator interface {
	Generate() string
}

type StateRepo interface {
	AddState(state string)
}
