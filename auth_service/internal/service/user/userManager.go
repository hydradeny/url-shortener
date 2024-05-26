package user

import (
	"context"

	"golang.org/x/exp/slog"
)

type UserStorage interface {
	Create(ctx context.Context, in *CreateUser) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type UserManager struct {
	log     slog.Logger
	storage UserStorage
}

func Create(ctx context.Context, in *CreateUser) (*User, error) {
	return nil, nil
}

func GetByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}

func CheckPasswordByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}
