package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mizanmahi/aiusage/server/internal/domain"
	"github.com/mizanmahi/aiusage/types"
)

type contextKey string

const userContextKey contextKey = "user"

type UserFinder interface {
	FindByAPIKeyHash(ctx context.Context, apiKeyHash string) (*domain.User, error)
}

func Auth(users UserFinder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := bearerToken(r.Header.Get("Authorization"))
			if apiKey == "" {
				writeAuthError(w)
				return
			}

			user, err := users.FindByAPIKeyHash(r.Context(), hashAPIKey(apiKey))
			if err != nil || user == nil {
				writeAuthError(w)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func User(r *http.Request) *domain.User {
	user, _ := r.Context().Value(userContextKey).(*domain.User)
	return user
}

func bearerToken(header string) string {
	prefix := "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

func hashAPIKey(apiKey string) string {
	sum := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(sum[:])
}

func writeAuthError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(types.APIErrorResponse{
		Error: types.APIError{
			Code:    "UNAUTHORIZED",
			Message: "invalid API key",
		},
	})
}
