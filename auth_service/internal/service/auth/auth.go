package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydradeny/url-shortener/auth_service/internal/service/user"

	"github.com/hydradeny/url-shortener/auth_service/internal/service/session"
)

type SessionManager interface {
	Create(context.Context, *session.CreateSession) (*session.Session, error)
	Check(context.Context, *session.CheckSession) (*session.Session, error)
	Destroy(context.Context, *session.DestroySession) error
	DestroyAll(context.Context, *session.DestroyAllSession) (int, error)
}

type UserManager interface {
	Create(ctx context.Context, in *user.CreateUser) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	CheckPasswordByEmail(ctx context.Context, in *user.CheckPassword) (*user.User, error)
}

type AuthService struct {
	log *slog.Logger
	sm  SessionManager
	um  UserManager
}

func NewService(log *slog.Logger, sm SessionManager, um UserManager) *AuthService {
	return &AuthService{
		log: log,
		sm:  sm,
		um:  um,
	}
}

func (s *AuthService) Login(ctx context.Context, in *LoginInput) (*LoginOutput, error) {
	const op = "AuthService.Login"
	err := in.Validate()
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		return nil, err
	}
	checkIn := user.CheckPassword{
		Email:    in.Email,
		Password: in.Password,
	}
	u, err := s.um.CheckPasswordByEmail(ctx, &checkIn)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		return nil, err
	}
	CreateSession := &session.CreateSession{UserID: u.ID}
	sess, err := s.sm.Create(ctx, CreateSession)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		return nil, err
	}

	return &LoginOutput{SessionID: sess.ID}, nil
}

func (s *AuthService) Logout(ctx context.Context, in *LogoutInput) error {
	const op = "AuthService.Logout"
	DestroySession := &session.DestroySession{SessionID: in.SessionID}
	err := s.sm.Destroy(ctx, DestroySession)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		return err
	}

	return nil
}

func (s *AuthService) Register(ctx context.Context, in *RegisterInput) (*RegisterOutput, error) {
	const op = "AuthService.Register"
	err := in.Validate()
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		return nil, err
	}

	CreateUser := user.CreateUser{
		Email:    in.Email,
		Password: in.Password,
	}
	u, err := s.um.Create(ctx, &CreateUser)
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		return nil, err
	}

	return &RegisterOutput{UserID: u.ID}, nil
}
