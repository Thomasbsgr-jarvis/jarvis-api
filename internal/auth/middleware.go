package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Thomasbsgr/jarvis-api/internal/web"
)

type contextKey string

const userIDKey contextKey = "userID"

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token manquant."})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := h.service.VerifyToken(r.Context(), tokenStr)
		if err != nil {
			if errors.Is(err, ErrInvalidToken) {
				web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token invalide."})
				return
			}
			slog.Error("AuthMiddleware: internal error during validation", "err", err)
			web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token invalide."})
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}
