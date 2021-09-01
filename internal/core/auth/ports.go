package auth

type Oauth2Provider interface {
	AuthURL() string
}
