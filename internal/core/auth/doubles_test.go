package auth_test

type Oauth2ProviderSpy struct {
	ReturnAuthURL string
}

func NewOauth2ProviderSpy() *Oauth2ProviderSpy {
	return &Oauth2ProviderSpy{ReturnAuthURL: "https://example.com/oauth"}
}

func (p *Oauth2ProviderSpy) AuthURL() string {
	return p.ReturnAuthURL
}
