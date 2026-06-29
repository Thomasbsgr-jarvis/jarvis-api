package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateUser
func (r *Repository) CreateUser(ctx context.Context, fullName, email, hashedPassword string) (*User, error) {
	var user User
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO users (full_name, email, password_hash) VALUES ($1, $2, $3) RETURNING id, full_name, email, password_hash, created_at",
		fullName, email, hashedPassword,
	).Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("DB: CreateUser: %w", err)
	}
	return &user, nil
}

// FindUserById
func (r *Repository) FindUserById(ctx context.Context, userId int64) (*User, error) {
	var user User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, full_name, email, password_hash, created_at FROM users WHERE id = $1",
		userId,
	).Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("DB: FindUserById: %w", err)
	}
	return &user, nil
}

// FindUserByEmail
func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, full_name, email, password_hash, created_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("DB: FindUserByEmail: %w", err)
	}
	return &user, nil
}

// CreateRefreshToken
func (r *Repository) CreateRefreshToken(ctx context.Context, userID int64, hashedToken string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID, hashedToken, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("DB: CreateRefreshToken: %w", err)
	}
	return nil
}

// DeleteRefreshToken
func (r *Repository) DeleteRefreshToken(ctx context.Context, hashedToken string) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM refresh_tokens WHERE token = $1",
		hashedToken,
	)
	if err != nil {
		return fmt.Errorf("DB: DeleteRefreshToken: %w", err)
	}
	return nil
}

// CleanUpRefreshToken
func (r *Repository) CleanUpRefreshToken(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM refresh_tokens WHERE user_id = $1 AND expires_at < NOW()",
		userID,
	)
	if err != nil {
		return fmt.Errorf("DB: CleanUpRefreshToken: %w", err)
	}
	return nil
}
