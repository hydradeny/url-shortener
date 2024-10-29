package authhandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/hydradeny/url-shortener/auth_service/internal/apperror"
	"github.com/hydradeny/url-shortener/auth_service/internal/service/auth"
	"github.com/hydradeny/url-shortener/auth_service/internal/service/session"
	"github.com/hydradeny/url-shortener/auth_service/pkg/restapi"
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

func New(log *slog.Logger, service AuthService) *AuthHandler {
	return &AuthHandler{
		log:     log,
		service: service,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" {
		h.log.Warn("HTTP login", slog.String("error", "email is required"))
		restapi.RespJSONError(w, http.StatusBadRequest, fmt.Errorf("email is required"))
		return
	}
	if password == "" {
		h.log.Warn("HTTP login", slog.String("error", "password is required"))
		restapi.RespJSONError(w, http.StatusBadRequest, fmt.Errorf("password is required"))
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
			switch appError.UserError {
			case apperror.ErrNotFound:
				fallthrough
			case apperror.ErrBadPassword:
				h.log.Warn("HTTP login", slog.String("error", err.Error()))
				restapi.RespJSONError(w, http.StatusUnauthorized, appError)
			case apperror.ErrInternal:
				h.log.Error("HTTP login", slog.String("error", err.Error()))
				restapi.RespJSONError(w, http.StatusInternalServerError, appError)
			default:
				h.log.Error("HTTP login", slog.String("error", err.Error()))
				restapi.RespJSONError(w, http.StatusInternalServerError, appError)
			}
			return
		}

		h.log.Error("HTTP login", slog.String("error", err.Error()))
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
		h.log.Warn("HTTP register", slog.String("error", "email is required"))
		restapi.RespJSONError(w, http.StatusBadRequest, fmt.Errorf("email is required"))
		return
	}
	if password == "" {
		h.log.Warn("HTTP register", slog.String("error", "password is required"))
		restapi.RespJSONError(w, http.StatusBadRequest, fmt.Errorf("password is required"))
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
			switch appError.UserError {
			case apperror.ErrUserExist:
				h.log.Warn("HTTP register", slog.String("error", err.Error()))
				restapi.RespJSONError(w, http.StatusConflict, appError)
			case apperror.ErrPasswordNotValid:
				h.log.Warn("HTTP register", slog.String("error", err.Error()))
				restapi.RespJSONError(w, http.StatusBadRequest, appError)
			case apperror.ErrEmailNotValid:
				h.log.Warn("HTTP register", slog.String("error", err.Error()))
				restapi.RespJSONError(w, http.StatusBadRequest, appError)
			case apperror.ErrInternal:
				h.log.Error("HTTP register", slog.String("error", err.Error()))
				restapi.RespJSONError(w, http.StatusInternalServerError, appError)
			default:
				h.log.Error("HTTP register", slog.String("error", err.Error()))
				restapi.RespJSONError(w, http.StatusInternalServerError, appError)
			}
			return
		}
		h.log.Error("HTTP register", slog.String("error", err.Error()))
		restapi.RespJSONError(w, http.StatusInternalServerError, apperror.ErrUnknown)
		return
	}

	restapi.RespJSON(w, res)
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := sessionFromContext(r.Context())
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
			switch appError.UserError {
			case apperror.ErrSessionNotFound:
				h.log.Warn("HTTP logout", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusUnauthorized)
			default:
				h.log.Error("HTTP logout", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		h.log.Error("HTTP logout", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func sessionFromContext(ctx context.Context) (*session.Session, error) {
	sess, ok := ctx.Value(CtxSessionKey).(*session.Session)
	if !ok {
		return nil, apperror.ErrNoAuth
	}
	return sess, nil
}
