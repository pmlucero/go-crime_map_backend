package security

import (
	"context"
	"net/http"
)

const (
	APIKeyHeader = "X-API-Key"
)

// contextKey es un tipo personalizado para las claves del contexto
type contextKey string

const (
	// UserIDKey es la clave para almacenar el ID del usuario en el contexto
	UserIDKey contextKey = "user_id"
)

// APIKeyMiddleware es un middleware que valida la API Key en las peticiones
func APIKeyMiddleware(repo *Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get(APIKeyHeader)
			if apiKey == "" {
				http.Error(w, "API Key requerida", http.StatusUnauthorized)
				return
			}

			// Validar la API Key contra la base de datos
			storedKey, err := repo.ValidateAPIKey(r.Context(), apiKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if storedKey == nil {
				http.Error(w, "API Key inválida", http.StatusUnauthorized)
				return
			}

			if storedKey.Status != "active" {
				http.Error(w, "API Key revocada", http.StatusUnauthorized)
				return
			}

			if storedKey.IsExpired() {
				http.Error(w, "API Key expirada", http.StatusUnauthorized)
				return
			}

			// Agregar el ID al contexto para uso posterior
			ctx := context.WithValue(r.Context(), UserIDKey, storedKey.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
