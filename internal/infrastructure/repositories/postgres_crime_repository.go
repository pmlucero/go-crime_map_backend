package repositories

import (
	"context"
	"fmt"
	"time"

	"go-crime_map_backend/internal/domain/entities"

	"github.com/jmoiron/sqlx"
)

// PostgresCrimeRepository implementa el repositorio de delitos usando PostgreSQL
type PostgresCrimeRepository struct {
	db *sqlx.DB
}

// NewPostgresCrimeRepository crea una nueva instancia del repositorio
func NewPostgresCrimeRepository(db *sqlx.DB) *PostgresCrimeRepository {
	return &PostgresCrimeRepository{
		db: db,
	}
}

// Create crea un nuevo delito
func (r *PostgresCrimeRepository) Create(ctx context.Context, crime *entities.Crime) error {
	query := `
		INSERT INTO crimes (
			id, title, description, crime_type, status,
			latitude, longitude, address,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10
		)
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
		return fmt.Errorf("error al crear el delito: %w", err)
	}

	return nil
}

// List obtiene una lista paginada de delitos
func (r *PostgresCrimeRepository) List(ctx context.Context, page, limit int) ([]entities.Crime, int64, error) {
	offset := (page - 1) * limit

	// Obtener total de registros
	var total int64
	err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM crimes WHERE deleted_at IS NULL")
	if err != nil {
		return nil, 0, fmt.Errorf("error al obtener total de delitos: %w", err)
	}

	// Obtener delitos
	query := `
		SELECT 
			id, title, description, crime_type, status,
			latitude, longitude, address,
			created_at, updated_at, deleted_at
		FROM crimes
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var crimes []entities.Crime
	err = r.db.SelectContext(ctx, &crimes, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error al obtener delitos: %w", err)
	}

	return crimes, total, nil
}

// GetByID obtiene un delito por su ID
func (r *PostgresCrimeRepository) GetByID(ctx context.Context, id string) (*entities.Crime, error) {
	query := `
		SELECT 
			id, title, description, crime_type, status,
			latitude, longitude, address,
			created_at, updated_at, deleted_at
		FROM crimes
		WHERE id = $1 AND deleted_at IS NULL
	`

	var crime entities.Crime
	err := r.db.GetContext(ctx, &crime, query, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el delito: %w", err)
	}

	return &crime, nil
}

// Update actualiza un delito existente
func (r *PostgresCrimeRepository) Update(ctx context.Context, crime *entities.Crime) error {
	query := `
		UPDATE crimes
		SET 
			title = $1,
			description = $2,
			crime_type = $3,
			status = $4,
			latitude = $5,
			longitude = $6,
			address = $7,
			updated_at = $8
		WHERE id = $9 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		crime.Title,
		crime.Description,
		crime.Type,
		crime.Status,
		crime.Location.Latitude,
		crime.Location.Longitude,
		crime.Location.Address,
		time.Now(),
		crime.ID,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar el delito: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("delito no encontrado")
	}

	return nil
}

// Delete realiza una eliminación lógica de un delito
func (r *PostgresCrimeRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE crimes
		SET 
			status = 'DELETED',
			deleted_at = $1,
			updated_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error al eliminar el delito: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("delito no encontrado")
	}

	return nil
}

// GetStats obtiene estadísticas de delitos
func (r *PostgresCrimeRepository) GetStats(ctx context.Context) (*entities.CrimeStats, error) {
	// Obtener total de delitos
	var total int64
	err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM crimes WHERE deleted_at IS NULL")
	if err != nil {
		return nil, fmt.Errorf("error al obtener total de delitos: %w", err)
	}

	// Obtener delitos por tipo
	var crimesByType map[string]int64
	err = r.db.SelectContext(ctx, &crimesByType, `
		SELECT crime_type, COUNT(*) as count
		FROM crimes
		WHERE deleted_at IS NULL
		GROUP BY crime_type
	`)
	if err != nil {
		return nil, fmt.Errorf("error al obtener delitos por tipo: %w", err)
	}

	// Obtener delitos por estado
	var crimesByStatus map[string]int64
	err = r.db.SelectContext(ctx, &crimesByStatus, `
		SELECT status, COUNT(*) as count
		FROM crimes
		WHERE deleted_at IS NULL
		GROUP BY status
	`)
	if err != nil {
		return nil, fmt.Errorf("error al obtener delitos por estado: %w", err)
	}

	// Obtener delitos por dirección
	var crimesByAddress map[string]int64
	err = r.db.SelectContext(ctx, &crimesByAddress, `
		SELECT address, COUNT(*) as count
		FROM crimes
		WHERE deleted_at IS NULL
		GROUP BY address
	`)
	if err != nil {
		return nil, fmt.Errorf("error al obtener delitos por dirección: %w", err)
	}

	return &entities.CrimeStats{
		TotalCrimes:     total,
		CrimesByType:    crimesByType,
		CrimesByStatus:  crimesByStatus,
		CrimesByAddress: crimesByAddress,
		LastUpdate:      time.Now(),
	}, nil
}
