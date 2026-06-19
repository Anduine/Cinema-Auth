package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	UserAgent string
	ClientIP  string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
