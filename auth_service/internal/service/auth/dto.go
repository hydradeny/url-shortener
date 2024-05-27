package auth

const (
	minPassLength = 8
	maxPassLength = 32
)

type LoginInput struct {
	Email    string
	Password string
}

func (d *LoginInput) Validate() error {
	// TODO: add validation
	return nil
}

type LogoutInput struct {
	SessionID string
}

type RegisterInput struct {
	Email    string
	Password string
}

func (d *RegisterInput) Validate() error {
	if len(d.Password) < minPassLength || len(d.Password) > maxPassLength {
		return ErrWrongPassLength
	}

	// TODO:Email valdation

	return nil
}

type RegisterOutput struct {
	UserID uint32
}

type LogoutOutput struct{}

type LoginOutput struct {
	SessionID string
}
