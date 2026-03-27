package middleware

import (
	"net/http"
)

func RequireScope(scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey, ok := GetApiKey(r.Context())
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"ok":false,"message":"authentication required"}`)) //nolint:errcheck
				return
			}

			if !hasScope(apiKey.Scopes, scope) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"ok":false,"message":"insufficient scope"}`)) //nolint:errcheck
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func hasScope(scopes []string, required string) bool {
	for _, s := range scopes {
		if s == "*" || s == required {
			return true
		}
	}
	return false
}
