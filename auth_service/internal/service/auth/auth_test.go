package auth_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/hydradeny/url-shortener/auth_service/internal/apperror"
	"github.com/hydradeny/url-shortener/auth_service/internal/service/auth"
	"github.com/hydradeny/url-shortener/auth_service/internal/service/session"
	"github.com/hydradeny/url-shortener/auth_service/internal/service/user"
	"github.com/hydradeny/url-shortener/auth_service/mocks"
)

func TestLoginCorrect(t *testing.T) {
	tc := []struct {
		name     string
		userID   uint32
		intput   auth.LoginInput
		expected auth.LoginOutput
	}{
		{
			"correct data",
			13,
			auth.LoginInput{
				Email:    "Vladimir@mail.Ru",
				Password: "1234567890",
			},

			auth.LoginOutput{
				SessionID: "id-1-2-3-5-6-7",
			},
		},
	}
	SMmock := mocks.NewSessionManager(t)
	USmock := mocks.NewUserManager(t)

	for i := range tc {

		USmock.EXPECT().CheckPasswordByEmail(context.Background(),
			&user.CheckPassword{
				Email:    tc[i].intput.Email,
				Password: tc[i].intput.Password,
			},
		).Return(
			&user.User{
				ID:    tc[i].userID,
				Email: tc[i].intput.Email,
			},
			nil,
		).Once()

		SMmock.EXPECT().Create(context.Background(),
			&session.CreateSession{
				UserID: tc[i].userID,
			},
		).Return(
			&session.Session{
				ID:     tc[i].expected.SessionID,
				UserID: tc[i].userID,
			},
			nil,
		).Once()

		auth := auth.NewService(nil, SMmock, USmock)
		got, err := auth.Login(context.Background(), &tc[i].intput)
		if err != nil {
			t.Errorf("expected nil error, got %s", err.Error())
		}
		if !reflect.DeepEqual(got, &tc[i].expected) {
			t.Errorf("expected %+v, got %+v", tc[i].expected, got)
		}

	}
}

func TestLoginIncorrectInput(t *testing.T) {
	tc := []struct {
		name      string
		intput    auth.LoginInput
		ErrNotNil error
	}{
		{
			"UserServiceErr",
			auth.LoginInput{
				Email:    "Vladimir@mail.Ru",
				Password: "1234567890",
			},
			apperror.ErrBadLogin,
		},
	}
	SMmock := mocks.NewSessionManager(t)
	USmock := mocks.NewUserManager(t)
	for i := range tc {

		USmock.EXPECT().CheckPasswordByEmail(context.Background(),
			&user.CheckPassword{
				Email:    tc[i].intput.Email,
				Password: tc[i].intput.Password,
			},
		).Return(
			nil,
			tc[i].ErrNotNil,
		).Once()

		auth := auth.NewService(nil, SMmock, USmock)
		_, err := auth.Login(context.Background(), &tc[i].intput)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if !SMmock.AssertNotCalled(t, "Create") {
			t.Errorf("SessionManager.Create called but shouldnt")
		}
	}
}
