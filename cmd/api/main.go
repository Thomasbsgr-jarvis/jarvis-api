package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Thomasbsgr/jarvis-api/internal/auth"
	"github.com/Thomasbsgr/jarvis-api/internal/config"
	"github.com/Thomasbsgr/jarvis-api/internal/database"
	"github.com/Thomasbsgr/jarvis-api/internal/eggo"
	"github.com/Thomasbsgr/jarvis-api/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()
	logger.Setup(cfg.AppEnv)

	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		slog.Error("Database connection failed", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	// Auth
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)

	// Eggo
	eggoRepo := eggo.NewRepository(db)
	eggoService := eggo.NewService(eggoRepo)
	eggoHandler := eggo.NewHandler(eggoService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestSize(1 << 20))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
		r.Post("/logout", authHandler.Logout)
	})

	r.Group(func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)

		r.Get("/me", authHandler.Me)
		r.Route("/eggo", func(r chi.Router) {
			r.Post("/complaints", eggoHandler.Complaints)
			r.Post("/files", eggoHandler.NewFile)
			r.Get("/complaints/{id}", eggoHandler.GetComplaint)
			r.Get("/files/{id}", eggoHandler.GetFiles)
		})
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("serveur error", "err", err)
			os.Exit(1)
		}
	}()

	slog.Info("server started", "port", cfg.Port)

	<-stop
	slog.Info("Stopping server…")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Graceful shutdown failed", "err", err)
	}
}
