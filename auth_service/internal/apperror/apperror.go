package apperror

import (
	"encoding/json"
	"errors"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrBadRequest       = errors.New("bad request")
	ErrUserExist        = errors.New("user already exists")
	ErrInternal         = errors.New("internal server error")
	ErrUserNotFound     = errors.New("user not found")
	ErrBadPassword      = errors.New("bad password")
	ErrUnknown          = errors.New("unknown error")
	ErrBadLogin         = errors.New("bad email or password")
	ErrPasswordNotValid = errors.New("validation password error")
	ErrSessionNotFound  = errors.New("session not found")
	ErrEmailNotValid    = errors.New("validation email error")
	ErrNoAuth           = errors.New("no session found")
)

type AppError struct {
	WrappedErr error `json:"-"`
	UserError  error
	Message    string `json:",omitempty"`
}

func (e *AppError) Unwrap() error {
	return e.WrappedErr
}

func (e *AppError) Error() string {
	return e.WrappedErr.Error() //+ ":" + e.Message
}

func (e *AppError) Is(err error) bool {
	return errors.Is(err, e.UserError)
}

func (e *AppError) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bytes
}

func NewAppError(errorType error, messqge string, err error) *AppError {
	return &AppError{
		UserError:  errorType,
		Message:    messqge,
		WrappedErr: err,
	}
}
