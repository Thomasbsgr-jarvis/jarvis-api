package auth

import (
	"context"
)

type ServiceInterface interface {
	Register(ctx context.Context, fullName, email, password string) (*AuthResponse, error)
	Login(ctx context.Context, email, password string) (*AuthResponse, error)
	Refresh(ctx context.Context, token string) (*AuthResponse, error)
	Logout(ctx context.Context, token string) error
	Me(ctx context.Context, userId int64) (*User, error)
}
