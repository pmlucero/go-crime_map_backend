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
	DB *sql.DB
}

// NewPostgresCrimeRepository crea una nueva instancia del repositorio
func NewPostgresCrimeRepository(db *sql.DB) *PostgresCrimeRepository {
	return &PostgresCrimeRepository{
		DB: db,
	}
}

// GetByID obtiene un delito por su ID
func (r *PostgresCrimeRepository) GetByID(ctx context.Context, id string) (*entities.Crime, error) {
	query := `
		SELECT id, title, description, crime_type, status, latitude, longitude,
		 address, address_number, city, province, country, zip_code, created_at, updated_at, deleted_at
		FROM crimes
		WHERE id = $1 AND deleted_at IS NULL
	`

	crime := &entities.Crime{}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&crime.ID,
		&crime.Title,
		&crime.Description,
		&crime.Type,
		&crime.Status,
		&crime.Location.Latitude,
		&crime.Location.Longitude,
		&crime.Location.Address,
		&crime.Location.AddressNumber,
		&crime.Location.City,
		&crime.Location.Province,
		&crime.Location.Country,
		&crime.Location.ZipCode,
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
		INSERT INTO crimes (id, title, description, crime_type, status, latitude, longitude, address, address_number, city, province, country, zip_code, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.DB.ExecContext(ctx, query,
		crime.ID,
		crime.Title,
		crime.Description,
		crime.Type,
		crime.Status,
		crime.Location.Latitude,
		crime.Location.Longitude,
		crime.Location.Address,
		crime.Location.AddressNumber,
		crime.Location.City,
		crime.Location.Province,
		crime.Location.Country,
		crime.Location.ZipCode,
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
	query := `
		SELECT id, title, description, crime_type, status, latitude, longitude,
		 address, address_number, city, province, country, zip_code, created_at, updated_at, deleted_at
		FROM crimes
		WHERE deleted_at IS NULL
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener delitos: %w", err)
	}
	defer rows.Close()

	var crimes []*entities.Crime
	for rows.Next() {
		crime := &entities.Crime{}
		err := rows.Scan(
			&crime.ID,
			&crime.Title,
			&crime.Description,
			&crime.Type,
			&crime.Status,
			&crime.Location.Latitude,
			&crime.Location.Longitude,
			&crime.Location.Address,
			&crime.Location.AddressNumber,
			&crime.Location.City,
			&crime.Location.Province,
			&crime.Location.Country,
			&crime.Location.ZipCode,
			&crime.CreatedAt,
			&crime.UpdatedAt,
			&crime.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear delito: %w", err)
		}
		crimes = append(crimes, crime)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar delitos: %w", err)
	}

	return crimes, nil
}

// Update actualiza un delito existente
func (r *PostgresCrimeRepository) Update(ctx context.Context, crime *entities.Crime) error {
	query := `
		UPDATE crimes
		SET title = $1, description = $2, crime_type = $3, status = $4,
			latitude = $5, longitude = $6, address = $7, address_number = $8, 
			city = $9, province = $10, country = $11, zip_code = $12, updated_at = $13,
			deleted_at = $14
		WHERE id = $15 AND deleted_at IS NULL
	`

	result, err := r.DB.ExecContext(ctx, query,
		crime.Title,
		crime.Description,
		crime.Type,
		crime.Status,
		crime.Location.Latitude,
		crime.Location.Longitude,
		crime.Location.Address,
		crime.Location.AddressNumber,
		crime.Location.City,
		crime.Location.Province,
		crime.Location.Country,
		crime.Location.ZipCode,
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
	query := `
		UPDATE crimes
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar delito: %w", err)
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

// List obtiene una lista de delitos con los filtros especificados
func (r *PostgresCrimeRepository) List(ctx context.Context, filter repositories.ListCrimesFilter) ([]*entities.Crime, error) {
	query := `
		SELECT id, title, description, crime_type, status, latitude, longitude,
		 address, address_number, city, province, country, zip_code, created_at, updated_at, deleted_at
		FROM crimes
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}
	argCount := 1

	if filter.Type != "" {
		query += fmt.Sprintf(" AND crime_type = $%d", argCount)
		args = append(args, filter.Type)
		argCount++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filter.Status)
		argCount++
	}

	if !filter.StartDate.IsZero() {
		query += fmt.Sprintf(" AND created_at >= $%d", argCount)
		args = append(args, filter.StartDate)
		argCount++
	}

	if !filter.EndDate.IsZero() {
		query += fmt.Sprintf(" AND created_at <= $%d", argCount)
		args = append(args, filter.EndDate)
		argCount++
	}

	if filter.Latitude != 0 && filter.Longitude != 0 && filter.RadiusKm > 0 {
		query += fmt.Sprintf(" AND earth_distance(ll_to_earth($%d, $%d), ll_to_earth(latitude, longitude)) <= $%d * 1000",
			argCount, argCount+1, argCount+2)
		args = append(args, filter.Latitude, filter.Longitude, filter.RadiusKm)
		argCount += 3
	}

	// Agregar paginación
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, filter.Limit, (filter.Page-1)*filter.Limit)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al listar delitos: %w", err)
	}
	defer rows.Close()

	var crimes []*entities.Crime
	for rows.Next() {
		crime := &entities.Crime{}
		err := rows.Scan(
			&crime.ID,
			&crime.Title,
			&crime.Description,
			&crime.Type,
			&crime.Status,
			&crime.Location.Latitude,
			&crime.Location.Longitude,
			&crime.Location.Address,
			&crime.Location.AddressNumber,
			&crime.Location.City,
			&crime.Location.Province,
			&crime.Location.Country,
			&crime.Location.ZipCode,
			&crime.CreatedAt,
			&crime.UpdatedAt,
			&crime.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear delito: %w", err)
		}
		crimes = append(crimes, crime)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre los delitos: %w", err)
	}

	return crimes, nil
}
