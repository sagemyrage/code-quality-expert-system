package middleware

import (
	"errors"
	"net/http"

	"github.com/sagemyrage/code-quality-expert-system/internal/http/requestcontext"
	"github.com/sagemyrage/code-quality-expert-system/internal/http/session"
	"github.com/sagemyrage/code-quality-expert-system/internal/service"
)

func IdentifyUser(authService *service.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(session.CookieName)
			if errors.Is(err, http.ErrNoCookie) {
				next.ServeHTTP(w, r)
				return
			}
			if err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			userID, err := authService.AuthenticateSession(r.Context(), cookie.Value)
			if errors.Is(err, service.ErrUnauthenticated) {
				next.ServeHTTP(w, r)
				return
			}
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			ctx := requestcontext.WithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := requestcontext.UserID(r.Context())
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
