package midleware

import (
	"context"
	"net/http"

	"github.com/hydradeny/url-shortener/auth_service/internal/service/session"
	"github.com/hydradeny/url-shortener/auth_service/internal/transport/http/handlers/authhandler"
)

func AuthMiddleware(sm session.SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionCookie, err := r.Cookie(authhandler.CookieName)
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		checkIn := &session.CheckSession{
			SessionID: sessionCookie.Value,
		}

		sess, err := sm.Check(ctx, checkIn)
		if err != nil {
			http.Error(w, "No auth", http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(ctx, authhandler.CtxSessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
