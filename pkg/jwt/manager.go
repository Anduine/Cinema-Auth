package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Role      string `json:"role"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

type Manager struct {
	secretKey string
	tokenTTL  time.Duration
}

func NewManager(secretKey string, tokenTTL time.Duration) *Manager {
	return &Manager{secretKey: secretKey, tokenTTL: tokenTTL}
}

func (m *Manager) Generate(userID, role, sessionID string) (string, error) {
	claims := CustomClaims{
		Role:      role,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *Manager) Verify(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Перевіряємо метод підпису
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("неочікуваний метод підпису")
			}
			return []byte(m.secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("недійсний токен")
	}

	return claims, nil
}
