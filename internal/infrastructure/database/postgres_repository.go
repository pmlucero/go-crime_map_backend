package database

import (
	"context"
	"database/sql"
	"fmt"
	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
)

// PostgresCrimeRepository implementa la interfaz CrimeRepository para PostgreSQL
type PostgresCrimeRepository struct {
	db *sql.DB
}

// NewPostgresCrimeRepository crea una nueva instancia del repositorio
func NewPostgresCrimeRepository(db *sql.DB) *PostgresCrimeRepository {
	return &PostgresCrimeRepository{
		db: db,
	}
}

// GetByID obtiene un delito por su ID
func (r *PostgresCrimeRepository) GetByID(ctx context.Context, id string) (*entities.Crime, error) {
	query := `
		SELECT id, title, description, crime_type, status, latitude, longitude, address, created_at, updated_at, deleted_at
		FROM crimes
		WHERE id = $1 AND deleted_at IS NULL
	`

	crime := &entities.Crime{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&crime.ID,
		&crime.Title,
		&crime.Description,
		&crime.Type,
		&crime.Status,
		&crime.Location.Latitude,
		&crime.Location.Longitude,
		&crime.Location.Address,
		&crime.CreatedAt,
		&crime.UpdatedAt,
		&crime.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error al obtener delito: %w", err)
	}

	return crime, nil
}

// Create crea un nuevo delito
func (r *PostgresCrimeRepository) Create(ctx context.Context, crime *entities.Crime) error {
	query := `
		INSERT INTO crimes (id, title, description, crime_type, status, latitude, longitude, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		crime.ID,
		crime.Title,
		crime.Description,
		crime.Type,
		crime.Status,
		crime.Location.Latitude,
		crime.Location.Longitude,
		crime.Location.Address,
		crime.CreatedAt,
		crime.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error al crear delito: %w", err)
	}

	return nil
}

// GetAll obtiene todos los delitos
func (r *PostgresCrimeRepository) GetAll(ctx context.Context) ([]*entities.Crime, error) {
	// TODO: Implementar
	return nil, nil
}

// Update actualiza un delito existente
func (r *PostgresCrimeRepository) Update(ctx context.Context, crime *entities.Crime) error {
	query := `
		UPDATE crimes
		SET title = $1, description = $2, crime_type = $3, status = $4,
			latitude = $5, longitude = $6, address = $7, updated_at = $8,
			deleted_at = $9
		WHERE id = $10
	`

	result, err := r.db.ExecContext(ctx, query,
		crime.Title,
		crime.Description,
		crime.Type,
		crime.Status,
		crime.Location.Latitude,
		crime.Location.Longitude,
		crime.Location.Address,
		crime.UpdatedAt,
		crime.DeletedAt,
		crime.ID,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar delito: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("delito no encontrado")
	}

	return nil
}

// Delete elimina un delito por su ID
func (r *PostgresCrimeRepository) Delete(ctx context.Context, id string) error {
	// TODO: Implementar
	return nil
}

// List obtiene una lista de delitos con los filtros especificados
func (r *PostgresCrimeRepository) List(ctx context.Context, filter repositories.ListCrimesFilter) ([]*entities.Crime, error) {
	// TODO: Implementar
	return nil, nil
}
