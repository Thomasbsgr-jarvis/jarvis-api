package auth

import (
	"context"
	"time"
)

type RepositoryInterface interface {
	CreateUser(ctx context.Context, fullName, email, hashedPassword string) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	CreateRefreshToken(ctx context.Context, userID int64, hashedToken string, expiresAt time.Time) error
	DeleteRefreshToken(ctx context.Context, hashedToken string) error
	RotateRefreshToken(ctx context.Context, hashedToken string) (*User, error)
	CleanUpRefreshToken(ctx context.Context, userID int64) error
	FindUserById(ctx context.Context, userId int64) (*User, error)
}
