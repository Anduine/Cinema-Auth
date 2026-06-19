package server

import (
	"log/slog"
	"net/http"

	"auth-service/internal/delivery/rest"
	"auth-service/internal/lib/responseHTTP"
	"auth-service/pkg/jwt"
	"auth-service/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter(h *rest.AuthHandler, jwtManager *jwt.Manager) *chi.Mux {
	r := chi.NewMux()

	r.Route("/api/v1/auth", func(r chi.Router) {

		// r.Use(middleware.RequestID)
		// r.Use(middleware.RealIP)
		// r.Use(middleware.Logger)
		// r.Use(middleware.Recoverer)

		r.Group(func(public chi.Router) {
			public.Post("/register", h.Register)

			public.Post("/login", h.Login)

			// Оновлення access токена за допомогою refresh токена з куки
			public.Post("/refresh", h.Refresh)
		})

		r.Group(func(protected chi.Router) {
			protected.Use(middleware.AuthMiddleware(jwtManager))

			protected.Post("/logout", h.Logout)

			protected.Get("/me", h.GetMe)

			// Завершити всі інші сесії (крім поточної)
			protected.Post("/sessions/revoke-all", h.RevokeAllSessions)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Маршрут не знайдено:", "method", r.Method, "path", r.URL.Path)
		responseHTTP.JSONError(w, "Маршрут не знайдено", http.StatusNotFound)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Заборонений метод:", "method", r.Method, "path", r.URL.Path)
		responseHTTP.JSONError(w, "Заборонений метод", http.StatusMethodNotAllowed)
	})

	return r
}
