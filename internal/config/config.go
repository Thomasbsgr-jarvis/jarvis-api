package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using system env")
	}
	return &Config{
		DatabaseURL: mustGet("DATABASE_URL"),
		JWTSecret:   mustGet("JWT_SECRET"),
		Port:        getOrDefault("API_PORT", "8080"),
	}
}

func mustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		slog.Error("Missing environment variable", "key", key)
		os.Exit(1)
	}
	return val
}

func getOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
