package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Thomasbsgr/jarvis-api/internal/config"
	"github.com/Thomasbsgr/jarvis-api/internal/web"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

// Register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FullName string `json:"fullName" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Warn("Register: deserialization error", "err", err)
		web.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	if err := config.ValidateData(input); err != nil {
		if errors.Is(err, config.ErrFields) {
			web.WriteJSON(w, http.StatusUnprocessableEntity, errors.Unwrap(err).Error())
			return
		}
		slog.Error("Register: validator error", "err", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	authResponse, err := h.service.Register(r.Context(), input.FullName, strings.ToLower(input.Email), input.Password)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			web.WriteJSON(w, http.StatusConflict, map[string]string{"message": "Cette adresse email est déjà utilisée."})
			return
		}
		slog.Error("Register: unexpected error", "err", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	web.WriteJSON(w, http.StatusCreated, authResponse)
}

// Login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Warn("Login: deserialization error", "err", err)
		web.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	if err := config.ValidateData(input); err != nil {
		if errors.Is(err, config.ErrFields) {
			web.WriteJSON(w, http.StatusUnprocessableEntity, errors.Unwrap(err).Error())
			return
		}
		slog.Error("Login: validator error", "err", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	authResponse, err := h.service.Login(r.Context(), strings.ToLower(input.Email), input.Password)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) || errors.Is(err, ErrInvalidCredentials) {
			web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "L’e-mail et/ou le mot de passe ne sont pas corrects."})
			return
		}
		slog.Error("Login: unexpected error", "err", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	web.WriteJSON(w, http.StatusOK, authResponse)
}

// Refresh
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token invalide."})
		return
	}

	var input struct {
		Token string `json:"token" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Warn("Refresh: deserialization error", "err", err)
		web.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	if err := config.ValidateData(input); err != nil {
		if errors.Is(err, config.ErrFields) {
			web.WriteJSON(w, http.StatusUnprocessableEntity, errors.Unwrap(err).Error())
			return
		}
		slog.Error("Refresh: validator error", "err", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	authResponse, err := h.service.Refresh(r.Context(), input.Token, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Utilisateur inexistant."})
			return
		}
		if errors.Is(err, ErrRefreshTokenNotFound) {
			web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token invalide."})
			return
		}
		slog.Error("Refresh: unexpected error", "err", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	web.WriteJSON(w, http.StatusOK, authResponse)
}

// Logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Token string `json:"token" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Warn("Logout: deserialization error", "err", err)
		web.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	if err := config.ValidateData(input); err != nil {
		if errors.Is(err, config.ErrFields) {
			web.WriteJSON(w, http.StatusUnprocessableEntity, errors.Unwrap(err).Error())
			return
		}
		slog.Error("Logout: validator error", "err", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	if err := h.service.Logout(r.Context(), input.Token); err != nil {
		slog.Error("Logout: unexpected error", "err", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	web.WriteJSON(w, http.StatusOK, map[string]string{"message": "Déconnecté."})
}

// Me
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token invalide."})
		return
	}

	user, err := h.service.Me(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			web.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "L'utilisateur n'existe pas."})
			return
		}

		slog.Error("Me: unexpected error", "err", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	web.WriteJSON(w, http.StatusOK, user)
}
