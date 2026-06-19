package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"auth-service/pkg/jwt"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey    contextKey = "userID"
	UserRoleKey  contextKey = "role"
	SessionIDKey contextKey = "sessionID"
)

func AuthMiddleware(jwtManager *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				slog.Debug("Не авторизовано, відсутній токен:", "error", err.Error())
				http.Error(w, `{"message": "Не авторизовано, відсутній токен}`, http.StatusUnauthorized)
				return
			}

			claims, err := jwtManager.Verify(cookie.Value)
			if err != nil {
				slog.Debug("Недійсний або протермінований токен", "error", err.Error())
				http.Error(w, `{"message": "Недійсний або протермінований токен"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.Subject)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
			ctx = context.WithValue(ctx, SessionIDKey, claims.SessionID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(UserRoleKey).(string)
		if !ok || role != "admin" {
			http.Error(w, `{"message": "Доступ заборонено, потрібні права адміністратора"}`, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func UserIDFromRequest(r *http.Request) uuid.UUID {
	userID, ok := r.Context().Value(UserIDKey).(string)
	userUUID, err := uuid.Parse(userID)
	if !ok || err != nil {
		slog.Debug("Не авторизовано:", "userID", userID, "contextVal", r.Context().Value(UserIDKey))
		return uuid.Nil
	}

	return userUUID
}

func UserRoleFromRequest(r *http.Request) string {
	return r.Context().Value(UserRoleKey).(string)
}

func UserRoleFromContext(ctx context.Context) string {
	return ctx.Value(UserRoleKey).(string)
}
