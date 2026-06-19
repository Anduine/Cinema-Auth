package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"auth-service/internal/domain"

	"github.com/google/uuid"
)

type postgresAuthRepo struct {
	db *sql.DB
}

func NewPostgresAuthRepo(db *sql.DB) domain.IAuthRepository {
	return &postgresAuthRepo{db: db}
}

func (r *postgresAuthRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *postgresAuthRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *postgresAuthRepo) CreateUser(ctx context.Context, user *domain.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return err
	}

	queryUser := `
		INSERT INTO users (id, username, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (email) DO NOTHING
	`

	resUser, err := tx.ExecContext(ctx, queryUser, user.ID, user.Username,
		user.Email, user.PasswordHash, user.Role, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		return err
	}

	rowsAffected, _ := resUser.RowsAffected()

	if rowsAffected == 0 {
		return domain.ErrUserAlreadyExist
	}

	queryProfile := `
		INSERT INTO user_profiles (user_id, updated_at)
		VALUES ($1, $2)
	`

	_, err = tx.ExecContext(ctx, queryProfile, user.ID, time.Now().UTC())
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
