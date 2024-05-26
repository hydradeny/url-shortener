package auth

import (
	"context"

	"github.com/hydradeny/url-shortener/auth_service/internal/service/user"
	"golang.org/x/exp/slog"

	"github.com/hydradeny/url-shortener/auth_service/internal/service/session"
)

type SessionManager interface {
	Create(context.Context, *session.CreateSession) (*session.Session, error)
	Check(context.Context, *session.CheckSession) (*session.Session, error)
	Destroy(context.Context, *session.DestroySession) error
	DestroyAll(context.Context, *session.DestroyAllSession) error
}

type UserManager interface {
	Create(ctx context.Context, in *user.CreateUser) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
}

type AuthService struct {
	log *slog.Logger
	sm  SessionManager
	um  UserManager
}

func (s *AuthService) Login(ctx context.Context, in *LoginInput) (*LoginOutput, error) {
	return nil, nil
}

func (s *AuthService) Logout(ctx context.Context, in *LogoutInput) (*LogoutOutput, error) {
	return nil, nil
}

func (s *AuthService) Register(ctx context.Context, in *RegisterInput) (*RegisterOutput, error) {
	return nil, nil
}
