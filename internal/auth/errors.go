package auth

import "errors"

var (
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrUserNotFound         = errors.New("user doesn't exist")
	ErrInvalidCredentials   = errors.New("invalid credential")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)
