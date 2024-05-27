package auth

import "errors"

var (
	ErrWrongPassLength = errors.New("wrong password length")
	UnknownErr         = errors.New("unknown error")
)
