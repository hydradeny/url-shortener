package session

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/hydradeny/url-shortener/auth_service/internal/apperror"
	"golang.org/x/exp/slog"
)

type SessionRepo interface {
	Create(context.Context, *CreateSessionStorage) (*Session, error)
	GetByID(ctx context.Context, id string) (*Session, error)
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID uint32) error
}

type SessionManager struct {
	log     *slog.Logger
	storage SessionRepo
}

func (sm *SessionManager) Create(ctx context.Context, in *CreateSession) (*Session, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, apperror.NewAppError(apperror.ErrInternal, "generation UUID error", err)
	}
	createSession := &CreateSessionStorage{
		ID:     uuid.String(),
		UserID: in.UserID,
	}
	sess, err := sm.storage.Create(ctx, createSession)
	if err != nil {
		return nil, err
	}
	return sess, err
}

func (sm *SessionManager) Check(ctx context.Context, in *CheckSession) (*Session, error) {
	sess, err := sm.storage.GetByID(ctx, in.SessionID)
	if err != nil {
		return nil, err
	}
	return sess, err
}

func (sm *SessionManager) Destroy(ctx context.Context, in *DestroySession) error {
	err := sm.storage.Delete(ctx, in.SessionID)
	if err != nil {
		return err
	}
	return nil
}

func (sm *SessionManager) DestroyAll(ctx context.Context, in *DestroyAllSession) error {
	err := sm.storage.DeleteByUserID(ctx, in.UserID)
	if err != nil {
		return err
	}
	return nil
}
