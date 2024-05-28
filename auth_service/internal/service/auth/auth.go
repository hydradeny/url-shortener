package auth

import (
	"context"
	"errors"

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
	CheckPasswordByEmail(ctx context.Context, email string) (*user.User, error)
}

type AuthService struct {
	log *slog.Logger
	sm  SessionManager
	um  UserManager
}

func (s *AuthService) Login(ctx context.Context, in *LoginInput) (*LoginOutput, error) {
	err := in.Validate()
	if err != nil {
		s.log.Error("validation", slog.String("error", err.Error()), "input:", *in)
		return nil, err
	}
	u, err := s.um.CheckPasswordByEmail(ctx, in.Email)
	if err != nil {
		s.log.Error("check password", slog.String("error", err.Error()))
		if errors.Is(err, user.UserNotFoundErr) || errors.Is(err, user.ErrBadPassword) {
			return nil, err
		}
		return nil, UnknownErr
	}
	CreateSession := &session.CreateSession{UserID: u.ID}
	sess, err := s.sm.Create(ctx, CreateSession)
	if err != nil {
		return nil, session.CreateSessionErr
	}

	return &LoginOutput{SessionID: sess.ID}, nil
}

func (s *AuthService) Logout(ctx context.Context, in *LogoutInput) error {
	DestroySession := &session.DestroySession{SessionID: in.SessionID}
	err := s.sm.Destroy(ctx, DestroySession)
	if err != nil {
		s.log.Error("destroy", slog.String("error", err.Error()))
		if errors.Is(err, session.ErrSessionNotFound) {
			return err
		}
		return UnknownErr
	}

	return nil
}

func (s *AuthService) Register(ctx context.Context, in *RegisterInput) (*RegisterOutput, error) {
	err := in.Validate()
	if err != nil {
		s.log.Error("validation", slog.String("error", err.Error()), "input:", *in)
		return nil, err
	}

	CreateUser := user.CreateUser{
		Email:    in.Email,
		Password: in.Password,
	}
	u, err := s.um.Create(ctx, &CreateUser)
	if err != nil {
		s.log.Error("user creation", slog.String("error", err.Error()), "input:", *in)
		if errors.Is(err, user.UserAlreadyExistErr) {
			return nil, user.UserAlreadyExistErr
		}
		return nil, user.UnknownUserCreationErr
	}

	return &RegisterOutput{UserID: u.ID}, nil
}
