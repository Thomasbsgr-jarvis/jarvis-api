package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Thomasbsgr/jarvis-api/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("DB opening error: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("DB connection error: %w", err)
	}

	return db, nil
}
