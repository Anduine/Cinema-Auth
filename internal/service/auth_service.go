package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"auth-service/internal/domain"
	"auth-service/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       domain.IAuthRepository
	jwtManager *jwt.Manager
}

func NewAuthService(repo domain.IAuthRepository, jwtManager *jwt.Manager) *AuthService {
	return &AuthService{repo: repo, jwtManager: jwtManager}
}

func (s *AuthService) Register(ctx context.Context, input domain.RegisterInput, userAgent, clientIP string) (*domain.User, *domain.TokenPair, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	userID, _ := uuid.NewV7()

	user := &domain.User{
		ID:           userID,
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Role:         domain.RoleUser,
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	pair, err := s.createSession(ctx, user, userAgent, clientIP)

	return user, pair, err
}

func (s *AuthService) Login(ctx context.Context, input domain.LoginInput, userAgent, clientIP string) (*domain.User, *domain.TokenPair, error) {
	user, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("неправильний email або пароль: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, nil, fmt.Errorf("неправильний email або пароль: %v", err)
	}

	pair, err := s.createSession(ctx, user, userAgent, clientIP)

	return user, pair, err
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	hash := s.hashToken(refreshToken)
	session, err := s.repo.GetSessionByTokenHash(ctx, hash)
	if err != nil {
		return nil
	}
	return s.repo.DeleteSession(ctx, session.ID)
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string, userAgent, clientIP string) (*domain.TokenPair, error) {
	hash := s.hashToken(refreshToken)

	session, err := s.repo.GetSessionByTokenHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("недійсний refresh токен: %v", err)
	}

	if session.ExpiresAt.Before(time.Now()) {
		_ = s.repo.DeleteSession(ctx, session.ID)
		return nil, fmt.Errorf("термін дії refresh токена вичерпано: %v", err)
	}

	user, err := s.repo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("користувача не знайдено: %v", err)
	}

	_ = s.repo.DeleteSession(ctx, session.ID)

	return s.createSession(ctx, user, userAgent, clientIP)
}

func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *AuthService) RevokeOtherSessions(ctx context.Context, userID, currentSessionID uuid.UUID) error {
	err := s.repo.DeleteAllExceptSessions(ctx, userID, currentSessionID)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) createSession(ctx context.Context, user *domain.User, userAgent, clientIP string) (*domain.TokenPair, error) {
	sessionID, _ := uuid.NewV7()

	rawRefreshToken, _ := uuid.NewV7()

	session := &domain.Session{
		ID:        sessionID,
		UserID:    user.ID,
		TokenHash: s.hashToken(rawRefreshToken.String()),
		UserAgent: userAgent,
		ClientIP:  clientIP,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	accessToken, err := s.jwtManager.Generate(user.ID.String(), user.Role, sessionID.String())
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken.String(),
	}, nil
}

func (s *AuthService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
