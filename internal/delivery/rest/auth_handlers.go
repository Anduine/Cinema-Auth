package rest

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"auth-service/internal/domain"
	"auth-service/internal/lib/responseHTTP"
	"auth-service/pkg/middleware"

	"github.com/google/uuid"
)

type AuthHandler struct {
	service         domain.IAuthService
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthHandler(service domain.IAuthService, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *AuthHandler {
	return &AuthHandler{service: service, accessTokenTTL: accessTokenTTL, refreshTokenTTL: refreshTokenTTL}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input domain.RegisterInput

	err := responseHTTP.DecodeJSONBody(w, r, &input)
	if err != nil {
		slog.Debug("Register: Невірний формат:", "error", err.Error())
		responseHTTP.JSONError(w, "Невірний формат даних", http.StatusBadRequest)
		return
	}

	userAgent := r.UserAgent()
	clientIP := r.RemoteAddr

	user, tokens, err := h.service.Register(r.Context(), input, userAgent, clientIP)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExist) {
			slog.Debug("Register: Користувач вже існує:", "error", err.Error())
			responseHTTP.JSONError(w, "Користувач вже існує", http.StatusConflict)
			return
		}

		slog.Debug("Register: Помилка при реєстрації користувача:", "error", err.Error())
		responseHTTP.JSONError(w, "Помилка при реєстрації", http.StatusInternalServerError)
		return
	}

	resp := &domain.UserResponse{
		Message: "Успішна реєстрація",
		Data:    user,
	}

	h.setTokenCookies(w, tokens)
	responseHTTP.JSONResp(w, resp, http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input domain.LoginInput

	err := responseHTTP.DecodeJSONBody(w, r, &input)
	if err != nil {
		slog.Debug("Login: Невірний формат:", "error", err.Error())
		responseHTTP.JSONError(w, "Невірний формат даних", http.StatusBadRequest)
		return
	}

	userAgent := r.UserAgent()
	clientIP := r.RemoteAddr

	user, tokens, err := h.service.Login(r.Context(), input, userAgent, clientIP)
	if err != nil {
		slog.Debug("Login: Помилка при логіні користувача:", "error", err.Error())
		responseHTTP.JSONError(w, "Помилка при логіні користувача", http.StatusUnauthorized)
		return
	}

	resp := &domain.UserResponse{
		Message: "Успішний вхід",
		Data:    user,
	}

	h.setTokenCookies(w, tokens)
	responseHTTP.JSONResp(w, resp, http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		// Якщо кука є, видаляємо сесію з БД
		err = h.service.Logout(r.Context(), cookie.Value)
		if err != nil {
			slog.Debug("Logout: Помилка при видалені токену:", "error", err.Error())
		}
	}

	h.clearTokenCookies(w)
	responseHTTP.JSONRespMessage(w, "Успішний вихід", http.StatusOK)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromRequest(r)
	if userID == uuid.Nil {
		slog.Debug("GetMe: Неправильний ID:", "userID", userID)
		responseHTTP.JSONRespMessage(w, "Неправильний ID у токені", http.StatusNotFound)
		return
	}

	user, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		slog.Debug("GetMe: Користувача не знайдено:", "userID", userID, "error", err.Error())
		responseHTTP.JSONRespMessage(w, "Користувача не знайдено", http.StatusNotFound)
		return
	}

	resp := domain.UserResponse{
		Data: &domain.User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}

	responseHTTP.JSONResp(w, resp, http.StatusOK)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		slog.Debug("Refresh: Refresh токен відсутній:", "error", err.Error())
		responseHTTP.JSONError(w, "Refresh токен відсутній", http.StatusUnauthorized)
		return
	}

	userAgent := r.UserAgent()
	clientIP := r.RemoteAddr

	tokens, err := h.service.RefreshTokens(r.Context(), cookie.Value, userAgent, clientIP)
	if err != nil {
		h.clearTokenCookies(w)
		slog.Debug("Refresh: Помилка при рефреші токену:", "error", err.Error())
		responseHTTP.JSONError(w, "Помилка при рефреші токену", http.StatusUnauthorized)
		return
	}

	h.setTokenCookies(w, tokens)
	responseHTTP.JSONRespMessage(w, "Токени оновлено", http.StatusOK)
}

func (h *AuthHandler) RevokeAllSessions(w http.ResponseWriter, r *http.Request) {
	userID, ok1 := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	sessionID, ok2 := r.Context().Value(middleware.SessionIDKey).(uuid.UUID)

	if !ok1 || !ok2 {
		slog.Debug("RevokeAllSessions: Помилка авторизації:", "userID", userID, "sessionID", sessionID)
		responseHTTP.JSONRespMessage(w, "Помилка авторизації", http.StatusUnauthorized)
		return
	}

	err := h.service.RevokeOtherSessions(r.Context(), userID, sessionID)
	if err != nil {
		slog.Debug("RevokeAllSessions: Помилка при завершенні сеансів:", "userID", userID, "sessionID", sessionID)
		responseHTTP.JSONRespMessage(w, "Помилка при завершенні сеансів", http.StatusInternalServerError)
		return
	}

	responseHTTP.JSONRespMessage(w, "Усі інші сеанси успішно завершено", http.StatusOK)
}

func (h *AuthHandler) setTokenCookies(w http.ResponseWriter, tokens *domain.TokenPair) {
	// Access Token живе недовго
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		HttpOnly: true,
		Secure:   false, // TODO: true для HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   int(h.accessTokenTTL.Seconds()), // 15 хвилин
	})

	// Refresh Token живе довго
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/api/v1/auth/refresh",           // Обмежуємо шлях, щоб не ганяти токен
		MaxAge:   int(h.refreshTokenTTL.Seconds()), // 30 днів
	})
}

func (h *AuthHandler) clearTokenCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/api/auth",
		MaxAge:   -1,
	})
}
