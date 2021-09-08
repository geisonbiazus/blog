package auth

type ConfimOauth2UseCase struct{}

func NewConfimOauth2UseCase() *ConfimOauth2UseCase {
	return &ConfimOauth2UseCase{}
}

func (u *ConfimOauth2UseCase) Run(state, code string) error {
	return ErrInvalidState
}
