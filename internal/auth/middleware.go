package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Thomasbsgr/jarvis-api/internal/web"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token manquant."})
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			}, jwt.WithExpirationRequired())
			if err != nil || !token.Valid {
				web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token invalide."})
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				slog.Error("AuthMiddleware: failed to cast claims")
				web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token invalide."})
				return
			}

			sub, err := claims.GetSubject()
			if err != nil || sub == "" {
				slog.Error("AuthMiddleware: failed to get subject claim")
				web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token invalide."})
				return
			}

			userID, err := strconv.ParseInt(sub, 10, 64)
			if err != nil {
				slog.Error("AuthMiddleware: failed to parse subject as int64", "sub", sub)
				web.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token invalide."})
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}
