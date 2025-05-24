package security

import (
	"context"
	"database/sql"

	"go-crime_map_backend/internal/domain/entities"
)

// Repository maneja las operaciones de base de datos para las API Keys
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// ValidateAPIKey valida una API Key y retorna sus detalles si es válida
func (r *Repository) ValidateAPIKey(ctx context.Context, key string) (*entities.APIKey, error) {
	query := `
		SELECT key, user_id, is_active, created_at, expires_at
		FROM api_keys
		WHERE key = ? AND is_active = 1
	`

	apiKey := &entities.APIKey{}
	err := r.db.QueryRowContext(ctx, query, key).Scan(
		&apiKey.Key,
		&apiKey.ID,
		&apiKey.IsActiveBool,
		&apiKey.CreatedAt,
		&apiKey.ExpiresAt,
	)

	if err != nil {
		// Si la tabla no existe, no hay coincidencia o hay otro error, denegar acceso
		return nil, err
	}

	return apiKey, nil
}

// CreateAPIKey almacena una nueva API Key en la base de datos
func (r *Repository) CreateAPIKey(ctx context.Context, apiKey *entities.APIKey) error {
	query := `
		INSERT INTO api_keys (key, user_id, is_active, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		apiKey.Key,
		apiKey.ID,
		apiKey.IsActiveBool,
		apiKey.ExpiresAt,
		apiKey.CreatedAt,
	)

	return err
}

// GetAPIKey obtiene una API Key por su valor
func (r *Repository) GetAPIKey(ctx context.Context, key string) (*entities.APIKey, error) {
	query := `
		SELECT key, user_id, is_active, expires_at, created_at
		FROM api_keys
		WHERE key = ?
	`

	apiKey := &entities.APIKey{}
	err := r.db.QueryRowContext(ctx, query, key).Scan(
		&apiKey.Key,
		&apiKey.ID,
		&apiKey.IsActiveBool,
		&apiKey.ExpiresAt,
		&apiKey.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return apiKey, err
}

// ListAPIKeys obtiene todas las API Keys
func (r *Repository) ListAPIKeys(ctx context.Context) ([]*entities.APIKey, error) {
	query := `
		SELECT key, user_id, is_active, expires_at, created_at
		FROM api_keys
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var apiKeys []*entities.APIKey
	for rows.Next() {
		apiKey := &entities.APIKey{}
		err := rows.Scan(
			&apiKey.Key,
			&apiKey.ID,
			&apiKey.IsActiveBool,
			&apiKey.ExpiresAt,
			&apiKey.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys, nil
}

// RevokeAPIKey desactiva una API Key
func (r *Repository) RevokeAPIKey(ctx context.Context, key string) error {
	query := `
		UPDATE api_keys
		SET is_active = false
		WHERE key = ?
	`

	_, err := r.db.ExecContext(ctx, query, key)
	return err
}
