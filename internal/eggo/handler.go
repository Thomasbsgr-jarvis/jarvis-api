package eggo

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Thomasbsgr/jarvis-api/internal/auth"
	"github.com/Thomasbsgr/jarvis-api/internal/config"
	"github.com/Thomasbsgr/jarvis-api/internal/web"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

// Complaints
func (h *Handler) Complaints(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FolderId string `json:"folderId" validate:"required"`
		Content  string `json:"content" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Warn("Complaints: deserialization error", "err", err)
		web.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	if err := config.ValidateData(input); err != nil {
		if errors.Is(err, config.ErrFields) {
			web.WriteJSON(w, http.StatusUnprocessableEntity, errors.Unwrap(err).Error())
			return
		}
		slog.Error("Complaints: validator error", "err", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token invalide."})
		return
	}

	complaint, err := h.service.Complaints(r.Context(), userID, input.FolderId, input.Content)
	if err != nil {
		if errors.Is(err, ErrComplaintAlreadyExists) {
			web.WriteJSON(w, http.StatusConflict, map[string]string{"id": complaint.ID})
			return
		}
		slog.Error("Register: unexpected error", "err", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	web.WriteJSON(w, http.StatusCreated, map[string]string{"id": complaint.ID})
}

// NewFile
func (h *Handler) NewFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token invalide."})
		return
	}

	var input struct {
		ComplaintId string `json:"complaintId" validate:"required"`
		Hash        string `json:"hash" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Url         string `json:"url" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Warn("NewFile: deserialization error", "err", err)
		web.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	if err := config.ValidateData(input); err != nil {
		if errors.Is(err, config.ErrFields) {
			web.WriteJSON(w, http.StatusUnprocessableEntity, errors.Unwrap(err).Error())
			return
		}
		slog.Error("Complaints: validator error", "err", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	if err := h.service.NewFile(r.Context(), userID, input.ComplaintId, input.Hash, input.Name, input.Url); err != nil {
		slog.Error("NewFile: unexpected error", "err", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Une erreur est survenue."})
		return
	}
	web.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Fichier créé avec succès."})
}

// GetComplaint
func (h *Handler) GetComplaint(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token invalide."})
		return
	}
	complaint, err := h.service.GetComplaint(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, ErrComplaintNotFound) {
			web.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "Plainte introuvable."})
			return
		}
		slog.Error("GetComplaintById: unexpected error", "err", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	web.WriteJSON(w, http.StatusOK, complaint)
}

// GetFiles
func (h *Handler) GetFiles(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token invalide."})
		return
	}

	files, err := h.service.GetFiles(r.Context(), id, userID)
	if err != nil {
		slog.Error("GetComplaintById: unexpected error", "err", err)
		web.WriteJSON(w, http.StatusInternalServerError, map[string]string{"message": "Une erreur est survenue."})
		return
	}

	web.WriteJSON(w, http.StatusOK, files)
}
