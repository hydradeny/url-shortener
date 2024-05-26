package session

import "errors"

var (
	CreateSessionErr   = errors.New("can't create session")
	ErrSessionNotFound = errors.New("sesion not found")
)
