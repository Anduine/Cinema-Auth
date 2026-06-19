package repository

import (
	"context"
	"database/sql"
	"errors"

	"auth-service/internal/domain"

	"github.com/google/uuid"
)

func (r *postgresAuthRepo) GetSessionByTokenHash(ctx context.Context, hash string) (*domain.Session, error) {
	query := `
		SELECT id, user_id, token_hash, user_agent, client_ip, expires_at, created_at
		FROM sessions
		WHERE token_hash = $1
	`

	session := &domain.Session{}
	err := r.db.QueryRowContext(ctx, query, hash).Scan(
		&session.ID, &session.UserID, &session.TokenHash,
		&session.UserAgent, &session.ClientIP,
		&session.ExpiresAt, &session.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("сесію не знайдено")
		}
		return nil, err
	}

	return session, nil
}

func (r *postgresAuthRepo) CreateSession(ctx context.Context, session *domain.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, token_hash, user_agent, client_ip, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx, query,
		session.ID, session.UserID, session.TokenHash,
		session.UserAgent, session.ClientIP, session.ExpiresAt,
	)

	return err
}

func (r *postgresAuthRepo) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	query := `
		DELETE FROM sessions
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, sessionID)
	return err
}

func (r *postgresAuthRepo) DeleteAllExceptSessions(ctx context.Context, userID uuid.UUID, currentSessionID uuid.UUID) error {
	query := `
		DELETE FROM sessions
		WHERE user_id = $1 AND id != $2
	`

	_, err := r.db.ExecContext(ctx, query, userID, currentSessionID)
	return err
}
