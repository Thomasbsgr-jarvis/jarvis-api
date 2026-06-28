package auth

import (
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	FullName     string    `json:"fullName"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
