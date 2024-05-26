package auth

import (
	"errors"

	"github.com/hydradeny/url-shortener/auth_service/internal/service/session"
)

const (
	minPassLength = 8
	maxPassLength = 32
)

type LoginInput struct {
	Email    string
	Password string
}

type LogoutInput struct {
	session session.Session
}

type RegisterInput struct {
	Email    string
	Password string
}

func (d *RegisterInput) Validate() error {
	if len(d.Password) < minPassLength || len(d.Password) > maxPassLength {
		return errors.New("wrong password length")
	}

	// TODO:Email valdation

	return nil
}

type RegisterOutput struct {
	UserID uint32
}

type LogoutOutput struct{}

type LoginOutput struct {
	session session.Session
}
