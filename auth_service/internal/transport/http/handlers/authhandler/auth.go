package authhandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hydradeny/url-shortener/auth_service/internal/apperror"
	"github.com/hydradeny/url-shortener/auth_service/internal/service/auth"
	"github.com/hydradeny/url-shortener/auth_service/internal/service/session"
	"github.com/hydradeny/url-shortener/auth_service/pkg/restapi"
	"golang.org/x/exp/slog"
)

const CookieName = "session_id"

type CtxKey int

var CtxSessionKey CtxKey = 1

type AuthService interface {
	Register(ctx context.Context, in *auth.RegisterInput) (*auth.RegisterOutput, error)
	Logout(ctx context.Context, in *auth.LogoutInput) error
	Login(ctx context.Context, in *auth.LoginInput) (*auth.LoginOutput, error)
}
type AuthHandler struct {
	log     *slog.Logger
	service AuthService
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" {
		restapi.RespJSONError(w, http.StatusBadRequest, fmt.Errorf("empty email"))
		return
	}
	if password == "" {
		restapi.RespJSONError(w, http.StatusBadRequest, fmt.Errorf("empty password"))
		return
	}
	loginInput := &auth.LoginInput{
		Email:    email,
		Password: password,
	}
	appError := &apperror.AppError{}
	res, err := h.service.Login(r.Context(), loginInput)
	if err != nil {
		if errors.As(err, &appError) {
			switch appError.ErrType {
			case apperror.ErrNotFound:
				fallthrough
			case apperror.ErrBadPassword:
				restapi.RespJSONError(w, http.StatusUnauthorized, apperror.ErrBadLogin)
			case apperror.ErrInternal:
				restapi.RespJSONError(w, http.StatusInternalServerError, apperror.ErrInternal)
			default:
				restapi.RespJSONError(w, http.StatusInternalServerError, apperror.ErrUnknown)
			}
			return
		}

		restapi.RespJSONError(w, http.StatusInternalServerError, apperror.ErrUnknown)
		return
	}
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    res.SessionID,
		Expires:  time.Now().Add(90 * 24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" {
		restapi.RespJSONError(w, http.StatusBadRequest, fmt.Errorf("empty email"))
		return
	}
	if password == "" {
		restapi.RespJSONError(w, http.StatusBadRequest, fmt.Errorf("empty password"))
		return
	}
	registerIn := &auth.RegisterInput{
		Email:    email,
		Password: password,
	}
	appError := &apperror.AppError{}
	res, err := h.service.Register(r.Context(), registerIn)
	if err != nil {
		if errors.As(err, &appError) {
			switch appError.ErrType {
			case apperror.ErrUserExist:
				restapi.RespJSONError(w, http.StatusConflict, apperror.ErrUserExist)
			case apperror.ErrPasswordNotValid:
				restapi.RespJSONError(w, http.StatusBadRequest, apperror.ErrPasswordNotValid)
			case apperror.ErrEmailNotValid:
				restapi.RespJSONError(w, http.StatusBadRequest, apperror.ErrEmailNotValid)
			case apperror.ErrInternal:
				restapi.RespJSONError(w, http.StatusInternalServerError, apperror.ErrInternal)
			default:
				restapi.RespJSONError(w, http.StatusInternalServerError, apperror.ErrUnknown)
			}
			return
		}
		restapi.RespJSONError(w, http.StatusInternalServerError, apperror.ErrUnknown)
		return
	}

	restapi.RespJSON(w, res)
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: check session

	session, err := SessionFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	logoutIn := &auth.LogoutInput{
		SessionID: session.ID,
	}
	appError := &apperror.AppError{}
	err = h.service.Logout(r.Context(), logoutIn)
	if err != nil {
		if errors.As(err, &appError) {
			switch appError.ErrType {
			case apperror.ErrSessionNotFound:
				w.WriteHeader(http.StatusUnauthorized)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SessionFromContext(ctx context.Context) (*session.Session, error) {
	sess, ok := ctx.Value(CtxSessionKey).(*session.Session)
	if !ok {
		return nil, apperror.ErrNoAuth
	}
	return sess, nil
}
