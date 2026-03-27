package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/ConflictHQ/boilerworks-go-micro/internal/database/queries"
)

type contextKey string

const ApiKeyContextKey contextKey = "apiKey"

func ApiKeyAuth(q *queries.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"ok":false,"message":"missing X-API-Key header"}`)) //nolint:errcheck
				return
			}

			hash := sha256.Sum256([]byte(key))
			keyHash := hex.EncodeToString(hash[:])

			apiKey, err := q.GetApiKeyByHash(r.Context(), keyHash)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"ok":false,"message":"invalid API key"}`)) //nolint:errcheck
				return
			}

			// Update last used timestamp in background
			go func() {
				_ = q.UpdateLastUsed(context.Background(), apiKey.ID)
			}()

			ctx := context.WithValue(r.Context(), ApiKeyContextKey, apiKey)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetApiKey(ctx context.Context) (queries.ApiKey, bool) {
	key, ok := ctx.Value(ApiKeyContextKey).(queries.ApiKey)
	return key, ok
}
