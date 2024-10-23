package postgres

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"testing"

	"github.com/hydradeny/url-shortener/auth_service/internal/apperror"
	"github.com/hydradeny/url-shortener/auth_service/internal/service/user"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
)

func TestCreateGood(t *testing.T) {
	testCases := []struct {
		in  user.CreateUser
		out user.RawUser
	}{
		{
			in: user.CreateUser{
				Email:    "email",
				Password: "password",
			},
			out: user.RawUser{
				Email:    "email",
				PassHash: []byte("password"),
				ID:       1,
			},
		},
	}
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	store := NewPgxUserRepo(context.Background(), mock, &slog.Logger{})

	for _, tc := range testCases {
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(tc.in.Email, tc.in.Password).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(tc.out.ID))

		res, err := store.Create(context.Background(), &tc.in)
		if err != nil {
			t.Errorf("expected nil error, got: %s", err)
		}
		if !reflect.DeepEqual(tc.out, *res) {
			t.Errorf("got: %v\n expected: %v", *res, tc.out)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestCreateBad(t *testing.T) {
	testCases := []struct {
		name            string
		expectedMockErr error
		expectedFuncErr error
		in              *user.CreateUser
	}{
		{
			name:            "user exists",
			expectedMockErr: pgx.ErrNoRows,
			expectedFuncErr: apperror.NewAppError(apperror.ErrUserExist, "", nil),
			in: &user.CreateUser{
				Email:    "email",
				Password: "password",
			},
		},
		{
			name:            "other db error",
			expectedMockErr: errors.New("some other error"),
			expectedFuncErr: apperror.NewAppError(apperror.ErrInternal, "", errors.New("some other error")),
			in: &user.CreateUser{
				Email:    "email",
				Password: "password",
			},
		},
	}

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	store := NewPgxUserRepo(context.Background(), mock, &slog.Logger{})
	for _, tc := range testCases {
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(tc.in.Email, tc.in.Password).WillReturnError(tc.expectedMockErr)
		// WillReturnRows(pgxmock.NewRows([]string{"id"}))

		_, err := store.Create(context.Background(), tc.in)
		if !errors.Is(err, tc.expectedFuncErr) {
			t.Errorf("expected %v, got: %v", tc.expectedFuncErr, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetByEmailGood(t *testing.T) {
	testCases := []struct {
		in  string
		out user.RawUser
	}{
		{
			in: "email",
			out: user.RawUser{
				Email:    "email",
				PassHash: []byte("password"),
				ID:       1,
			},
		},
	}
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	store := NewPgxUserRepo(context.Background(), mock, &slog.Logger{})

	for _, tc := range testCases {
		mock.ExpectQuery("SELECT id, password FROM users").
			WithArgs(tc.in).
			WillReturnRows(pgxmock.NewRows([]string{"id", "password"}).AddRow(tc.out.ID, tc.out.PassHash))

		res, err := store.GetByEmail(context.Background(), tc.in)
		if err != nil {
			t.Errorf("expected nil error, got: %s", err)
		}
		if !reflect.DeepEqual(tc.out, *res) {
			t.Errorf("got: %v\n expected: %v", *res, tc.out)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetByEmailBad(t *testing.T) {
	testCases := []struct {
		name            string
		expectedMockErr error
		expectedFuncErr error
		in              string
	}{
		{
			name:            "user exists",
			expectedMockErr: pgx.ErrNoRows,
			expectedFuncErr: apperror.NewAppError(apperror.ErrNotFound, "", nil),
			in:              "email",
		},
		{
			name:            "other db error",
			expectedMockErr: errors.New("some other error"),
			expectedFuncErr: apperror.NewAppError(apperror.ErrInternal, "", errors.New("some other error")),
			in:              "email",
		},
	}

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	store := NewPgxUserRepo(context.Background(), mock, &slog.Logger{})
	for _, tc := range testCases {
		mock.ExpectQuery("SELECT id, password FROM users").
			WithArgs(tc.in).WillReturnError(tc.expectedMockErr)
		// WillReturnRows(pgxmock.NewRows([]string{"id"}))

		_, err := store.GetByEmail(context.Background(), tc.in)
		if !errors.Is(err, tc.expectedFuncErr) {
			t.Errorf("expected %v, got: %v", tc.expectedFuncErr, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}
