package apperror

import (
	"encoding/json"
	"errors"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)

type AppError struct {
	WrappedErr error
	ErrType    error
	Message    string
}

func (e *AppError) Unwrap() error {
	return e.WrappedErr
}

func (e *AppError) Error() string {
	return e.ErrType.Error() + ":" + e.Message + ":" + e.WrappedErr.Error()
}

func (e *AppError) Is(err error) bool {
	return errors.Is(err, e.ErrType)
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
		ErrType:    errorType,
		Message:    messqge,
		WrappedErr: err,
	}
}
