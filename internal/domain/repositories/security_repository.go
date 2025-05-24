package repositories

import (
	"context"

	"go-crime_map_backend/internal/domain/entities"
)

// SecurityRepository define las operaciones relacionadas con la seguridad
type SecurityRepository interface {
	// ValidateAPIKey valida una clave de API
	ValidateAPIKey(ctx context.Context, apiKey string) (*entities.APIKey, error)
}
