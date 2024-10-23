package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/hydradeny/url-shortener/auth_service/internal/apperror"
	"github.com/hydradeny/url-shortener/auth_service/internal/service/session"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgxPoolIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Ping(context.Context) error
	Close()
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

type PgxSessionRepo struct {
	dbpool PgxPoolIface
}

func NewPgxSessionRepo(pgxPool PgxPoolIface) *PgxSessionRepo {
	return &PgxSessionRepo{
		dbpool: pgxPool,
	}
}

func (repo *PgxSessionRepo) Create(ctx context.Context, in *session.CreateSessionStorage) error {
	const op = "PgxSessionRepo.Create"
	_, err := repo.dbpool.Exec(ctx, "INSERT INTO sessions (id, user_id) VALUES ($1,$2)", in.ID, in.UserID)
	if err != nil {
		return apperror.NewAppError(apperror.ErrInternal, "", fmt.Errorf("%s: %w", op, err))
	}
	return nil
}

func (repo *PgxSessionRepo) GetByID(ctx context.Context, id string) (*session.Session, error) {
	const op = "PgxSessionRepo.GetByID"
	row := repo.dbpool.QueryRow(ctx, "SELECT id, user_id FROM sessions WHERE id=$1", id)
	sess := session.Session{}
	err := row.Scan(&sess.ID, &sess.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewAppError(apperror.ErrNotFound, "", fmt.Errorf("%s: %w", op, err))
		}
		return nil, apperror.NewAppError(apperror.ErrInternal, "", fmt.Errorf("%s: %w", op, err))
	}
	return &sess, nil
}

func (repo *PgxSessionRepo) Delete(ctx context.Context, id string) error {
	const op = "PgxSessionRepo.Delete"
	_, err := repo.dbpool.Exec(ctx, "DELETE FROM sessions WHERE id=$1", id)
	if err != nil {
		return apperror.NewAppError(apperror.ErrInternal, "", fmt.Errorf("%s: %w", op, err))
	}
	return nil
}

func (repo *PgxSessionRepo) DeleteByUserID(ctx context.Context, userID uint32) (int, error) {
	const op = "PgxSessionRepo.DeleteByUserID"
	tag, err := repo.dbpool.Exec(ctx, "DELETE FROM sessions WHERE user_id=$1", userID)
	if err != nil {
		return 0, apperror.NewAppError(apperror.ErrInternal, "", fmt.Errorf("%s: %w", op, err))
	}
	return int(tag.RowsAffected()), nil
}
