package postgres

import (
	"context"
	"errors"

	"github.com/hydradeny/url-shortener/auth_service/internal/apperror"
	"github.com/hydradeny/url-shortener/auth_service/internal/service/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slog"
)

type PgxUserRepo struct {
	dbpool *pgxpool.Pool
	log    *slog.Logger
}

func NewPgxUserRepo(ctx context.Context, pgxPool *pgxpool.Pool, log *slog.Logger) *PgxUserRepo {
	return &PgxUserRepo{
		dbpool: pgxPool,
		log:    log,
	}
}

func (repo *PgxUserRepo) Create(ctx context.Context, in *user.CreateUser) (*user.RawUser, error) {
	row := repo.dbpool.QueryRow(ctx, "INSERT INTO users(email, password) values ($1,$2) RETURNING id", &in.Email, in.Password)
	out := &user.RawUser{
		PassHash: []byte(in.Password),
		Email:    in.Email,
	}
	err := row.Scan(&out.ID)
	if err != nil {
		// TODO: is it already exists error?
		if err == pgx.ErrNoRows {
			return nil, apperror.NewAppError(apperror.ErrUserExist, "", nil)
		}
		return nil, apperror.NewAppError(apperror.ErrInternal, "", err)
	}
	return out, nil
}

func (repo *PgxUserRepo) GetByEmail(ctx context.Context, email string) (*user.RawUser, error) {
	row := repo.dbpool.QueryRow(ctx, "SELECT id, password FROM users WHERE email=$1", email)
	out := &user.RawUser{
		Email: email,
	}
	err := row.Scan(&out.ID, &out.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewAppError(apperror.ErrNotFound, "email", nil)
		}
		return nil, apperror.NewAppError(apperror.ErrInternal, "", err)
	}
	return out, nil
}