package domain

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type IAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	CreateUser(ctx context.Context, user *User) error

	GetSessionByTokenHash(ctx context.Context, hash string) (*Session, error)
	CreateSession(ctx context.Context, session *Session) error
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	DeleteAllExceptSessions(ctx context.Context, userID uuid.UUID, currentSessionID uuid.UUID) error
}

type IAuthService interface {
	Register(ctx context.Context, input RegisterInput, userAgent, clientIP string) (*User, *TokenPair, error)

	Login(ctx context.Context, input LoginInput, userAgent, clientIP string) (*User, *TokenPair, error)

	RefreshTokens(ctx context.Context, refreshToken string, userAgent, clientIP string) (*TokenPair, error)

	Logout(ctx context.Context, refreshToken string) error

	GetProfile(ctx context.Context, userID uuid.UUID) (*User, error)

	RevokeOtherSessions(ctx context.Context, userID, currentSessionID uuid.UUID) error
}

type IAuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	GetMe(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	RevokeAllSessions(w http.ResponseWriter, r *http.Request)
}
