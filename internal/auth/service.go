package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      RepositoryInterface
	jwtSecret string
}

func NewService(repo RepositoryInterface, jwtSecret string) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret}
}

// Register
func (s *Service) Register(ctx context.Context, fullName, email, password string) (*AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("password hash error: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, fullName, email, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	authResponse, err := s.generateTokens(ctx, user)
	if err != nil {
		return nil, err
	}

	return authResponse, nil
}

// Login
func (s *Service) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := s.repo.CleanUpRefreshToken(ctx, user.ID); err != nil {
		return nil, err
	}

	authResponse, err := s.generateTokens(ctx, user)
	if err != nil {
		return nil, err
	}

	return authResponse, nil
}

// Refresh
func (s *Service) Refresh(ctx context.Context, token string) (*AuthResponse, error) {
	hashedToken := hashToken(token)

	user, err := s.repo.RotateRefreshToken(ctx, hashedToken)
	if err != nil {
		return nil, err
	}

	authResponse, err := s.generateTokens(ctx, user)
	if err != nil {
		return nil, err
	}

	return authResponse, nil
}

// Logout
func (s *Service) Logout(ctx context.Context, token string) error {
	hashedToken := hashToken(token)

	if err := s.repo.DeleteRefreshToken(ctx, hashedToken); err != nil {
		return err
	}

	return nil
}

// Me
func (s *Service) Me(ctx context.Context, userId int64) (*User, error) {
	user, err := s.repo.FindUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Token
func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("refresh token generation error: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func generateAccessToken(userID int64, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(userID, 10),
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error generating access token: %w", err)
	}
	return signed, nil
}

func (s *Service) generateTokens(ctx context.Context, user *User) (*AuthResponse, error) {
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	hashedRefreshToken := hashToken(refreshToken)

	if err := s.repo.CreateRefreshToken(ctx, user.ID, hashedRefreshToken, time.Now().Add(30*24*time.Hour)); err != nil {
		return nil, err
	}

	accessToken, err := generateAccessToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
