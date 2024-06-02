package session

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid"
	"github.com/hydradeny/url-shortener/auth_service/internal/apperror"
)

type SessionRepo interface {
	Create(context.Context, *CreateSessionStorage) error
	GetByID(ctx context.Context, id string) (*Session, error)
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID uint32) (deleted int, err error)
}

type SessionManager struct {
	log     *slog.Logger
	storage SessionRepo
}

func New(log *slog.Logger, repo SessionRepo) *SessionManager {
	return &SessionManager{
		log:     log,
		storage: repo,
	}
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
	err = sm.storage.Create(ctx, createSession)
	if err != nil {
		return nil, err
	}
	res := &Session{
		ID:     createSession.ID,
		UserID: createSession.UserID,
	}
	return res, err
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

func (sm *SessionManager) DestroyAll(ctx context.Context, in *DestroyAllSession) (int, error) {
	deleted, err := sm.storage.DeleteByUserID(ctx, in.UserID)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
