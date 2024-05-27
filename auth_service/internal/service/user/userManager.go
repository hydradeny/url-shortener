package user

import (
	"bytes"
	"context"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
	"golang.org/x/exp/slog"
)

const saltLength = 8

type UserStorage interface {
	Create(ctx context.Context, in *CreateUser) (*RawUser, error)
	GetByEmail(ctx context.Context, email string) (*RawUser, error)
}

type UserManager struct {
	log     slog.Logger
	storage UserStorage
}

func (um *UserManager) Create(ctx context.Context, in *CreateUser) (*User, error) {
	salt := makeSalt(saltLength)
	in.Password = string(um.hashPass(in.Password, salt))
	rawUser, err := um.storage.Create(ctx, in)
	if err != nil {
		return nil, err
	}
	user := &User{
		Email: rawUser.Email,
		ID:    rawUser.ID,
	}
	return user, nil
}

func (um *UserManager) GetByEmail(ctx context.Context, email string) (*User, error) {
	rawUser, err := um.storage.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	user := &User{
		Email: rawUser.Email,
		ID:    rawUser.ID,
	}
	return user, nil
}

func (um *UserManager) CheckPasswordByEmail(ctx context.Context, in *CheckPassword) (*User, error) {
	rawUser, err := um.storage.GetByEmail(ctx, in.Email)
	if err != nil {
		return nil, err
	}

	salt := rawUser.PassHash[0:saltLength]
	if !bytes.Equal(um.hashPass(in.Password, salt), rawUser.PassHash) {
		return nil, ErrBadPassword
	}

	user := &User{
		Email: rawUser.Email,
		ID:    rawUser.ID,
	}
	return user, nil
}

func makeSalt(n int) []byte {
	salt := make([]byte, n)
	rand.Read(salt)
	return salt
}

func (repo *UserManager) hashPass(plainPassword string, salt []byte) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)
	res := []byte(salt)
	return append(res, hashedPass...)
}
