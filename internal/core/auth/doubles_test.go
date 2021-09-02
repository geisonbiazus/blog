package auth_test

type Oauth2ProviderSpy struct {
	ReturnAuthURL string
	ReceivedState string
}

func NewOauth2ProviderSpy() *Oauth2ProviderSpy {
	return &Oauth2ProviderSpy{ReturnAuthURL: "https://example.com/oauth"}
}

func (p *Oauth2ProviderSpy) AuthURL(state string) string {
	p.ReceivedState = state
	return p.ReturnAuthURL
}

type IDGeneratorStub struct {
	ReturnID string
}

func NewIDGeneratorStub() *IDGeneratorStub {
	return &IDGeneratorStub{}
}

func (g *IDGeneratorStub) Generate() string {
	return g.ReturnID
}
