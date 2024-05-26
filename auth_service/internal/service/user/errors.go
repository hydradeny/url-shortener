package user

import "errors"

var (
	UserAlreadyExistErr    = errors.New("user already exists")
	UnknownUserCreationErr = errors.New("unknown user creation error")
	UserNotFoundErr        = errors.New("user not found")
	BadPasswordErr         = errors.New("bad password")
)
