package postgres

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/hydradeny/url-shortener/auth_service/internal/apperror"
	"github.com/hydradeny/url-shortener/auth_service/internal/service/user"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"golang.org/x/exp/slog"
)

func TestCreateGood(t *testing.T) {
	testCases := []struct {
		expectingMockErr error
		in               user.CreateUser
		out              user.RawUser
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
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(tc.out.ID)).RowsWillBeClosed()

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
			expectedMockErr: errors.New("some error"),
			expectedFuncErr: apperror.NewAppError(apperror.ErrInternal, "", errors.New("some error")),
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
