package user

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"

	"github.com/hydradeny/url-shortener/auth_service/internal/apperror"
	"golang.org/x/crypto/argon2"
)

const saltLength = 8

type UserStorage interface {
	Create(ctx context.Context, in *CreateUser) (*RawUser, error)
	GetByEmail(ctx context.Context, email string) (*RawUser, error)
}

type UserManager struct {
	log     *slog.Logger
	storage UserStorage
}

func NewService(log *slog.Logger, repo UserStorage) *UserManager {
	return &UserManager{
		log:     log,
		storage: repo,
	}
}

func (um *UserManager) Create(ctx context.Context, in *CreateUser) (*User, error) {
	const op = "UserManager.Create"
	salt := makeSalt(saltLength)
	in.Password = string(um.hashPass(in.Password, salt))
	rawUser, err := um.storage.Create(ctx, in)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		return nil, err
	}
	user := &User{
		Email: rawUser.Email,
		ID:    rawUser.ID,
	}
	return user, nil
}

func (um *UserManager) GetByEmail(ctx context.Context, email string) (*User, error) {
	const op = "UserManager.GetByEmail"
	rawUser, err := um.storage.GetByEmail(ctx, email)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		return nil, err
	}
	user := &User{
		Email: rawUser.Email,
		ID:    rawUser.ID,
	}
	return user, nil
}

func (um *UserManager) CheckPasswordByEmail(ctx context.Context, in *CheckPassword) (*User, error) {
	const op = "UserManager.CheckPasswordByEmail"
	rawUser, err := um.storage.GetByEmail(ctx, in.Email)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		return nil, err
	}

	salt := rawUser.PassHash[0:saltLength]
	if !bytes.Equal(um.hashPass(in.Password, salt), rawUser.PassHash) {
		return nil, apperror.NewAppError(apperror.ErrBadPassword, "", fmt.Errorf("%s: %w", op, err))
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
