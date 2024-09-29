package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hydradeny/url-shortener/auth_service/internal/service/auth"
	"github.com/hydradeny/url-shortener/auth_service/internal/service/session"
	"github.com/hydradeny/url-shortener/auth_service/internal/service/user"
	pgxsession "github.com/hydradeny/url-shortener/auth_service/internal/storage/session/postgres"
	pgxuser "github.com/hydradeny/url-shortener/auth_service/internal/storage/user/postgres"
	"github.com/hydradeny/url-shortener/auth_service/internal/transport/http/handlers/authhandler"
	mw "github.com/hydradeny/url-shortener/auth_service/internal/transport/http/midleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type config struct {
	env         string
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration
	pgConnStr   string
}

var cfg = config{
	pgConnStr:   "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable",
	env:         "local",
	Address:     "0.0.0.0:44044",
	Timeout:     time.Duration(time.Second * 4),
	IdleTimeout: time.Duration(time.Second * 30),
}

func main() {
	router := chi.NewRouter()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	dbpool, err := pgxpool.New(context.Background(), cfg.pgConnStr)
	if err != nil {
		log.Error("failed to create storage:", slog.String("error", err.Error()))
		return
	}
	userStorage := pgxuser.NewPgxUserRepo(context.Background(), dbpool, log)
	sessionStorage := pgxsession.NewPgxSessionRepo(dbpool)
	sm := session.New(log, sessionStorage)
	um := user.NewService(log, userStorage)
	authService := auth.NewService(log, sm, um)
	authHandler := authhandler.New(log, authService)

	router.Use(middleware.RequestID)
	// router.Use(middleware.Logger)

	router.Use(mw.NewLoggerMiddleware(log))
	router.Use(middleware.URLFormat)
	router.Use(middleware.Recoverer)

	router.Group(func(r chi.Router) {
		r.Use(mw.NewAuthMiddleware(sm))
		r.Get("/logout", authHandler.Logout)
	})

	router.Post("/login", authHandler.Login)
	router.Post("/reg", authHandler.Register)
	server := &http.Server{
		Addr:        cfg.Address,
		Handler:     router,
		ReadTimeout: cfg.Timeout,
		IdleTimeout: cfg.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error(err.Error())
	}
}
