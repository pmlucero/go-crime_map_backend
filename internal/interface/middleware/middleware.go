package middleware

import (
	"net/http"

	"go-crime_map_backend/internal/domain/repositories"
)

func APIKeyMiddleware(repo repositories.SecurityRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				http.Error(w, "API key is required", http.StatusUnauthorized)
				return
			}

			key, err := repo.ValidateAPIKey(r.Context(), apiKey)
			if err != nil {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			if !key.IsActiveBool {
				http.Error(w, "API key is not active", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
