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
			latitude, longitude, address, address_number, city, province, country, zip_code,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11, $12, $13,
			$14, $15
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
		crime.Location.AddressNumber,
		crime.Location.City,
		crime.Location.Province,
		crime.Location.Country,
		crime.Location.ZipCode,
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
			latitude as "location.latitude",
			longitude as "location.longitude",
			address as "location.address",
			address_number as "location.address_number",
			city as "location.city",
			province as "location.province",
			country as "location.country",
			zip_code as "location.zip_code",
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
			latitude as "location.latitude",
			longitude as "location.longitude",
			address as "location.address",
			address_number as "location.address_number",
			city as "location.city",
			province as "location.province",
			country as "location.country",
			zip_code as "location.zip_code",
			created_at, updated_at, deleted_at
		FROM crimes
		WHERE id = $1 AND deleted_at IS NULL
	`

	var crime entities.Crime
	err := r.db.GetContext(ctx, &crime, query, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el delito: %w", err)
	}

	// Mapear los campos de la ubicación
	crime.Location = entities.Location{
		Latitude:      crime.Location.Latitude,
		Longitude:     crime.Location.Longitude,
		Address:       crime.Location.Address,
		AddressNumber: crime.Location.AddressNumber,
		City:          crime.Location.City,
		Province:      crime.Location.Province,
		Country:       crime.Location.Country,
		ZipCode:       crime.Location.ZipCode,
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
			address_number = $8,
			city = $9,
			province = $10,
			country = $11,
			zip_code = $12,
			updated_at = $13
		WHERE id = $14 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
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

	// Obtener delitos activos
	var active int64
	err = r.db.GetContext(ctx, &active, "SELECT COUNT(*) FROM crimes WHERE deleted_at IS NULL AND status = 'ACTIVE'")
	if err != nil {
		return nil, fmt.Errorf("error al obtener delitos activos: %w", err)
	}

	// Obtener delitos inactivos
	var inactive int64
	err = r.db.GetContext(ctx, &inactive, "SELECT COUNT(*) FROM crimes WHERE deleted_at IS NULL AND status = 'INACTIVE'")
	if err != nil {
		return nil, fmt.Errorf("error al obtener delitos inactivos: %w", err)
	}

	// Obtener delitos por tipo
	type CrimeTypeCount struct {
		Type  string
		Count int64
	}
	var crimesByType []CrimeTypeCount
	err = r.db.SelectContext(ctx, &crimesByType, `
		SELECT crime_type as type, COUNT(*) as count
		FROM crimes
		WHERE deleted_at IS NULL
		GROUP BY crime_type
	`)
	if err != nil {
		return nil, fmt.Errorf("error al obtener delitos por tipo: %w", err)
	}

	// Convertir a map
	crimesByTypeMap := make(map[string]int64)
	for _, ct := range crimesByType {
		crimesByTypeMap[ct.Type] = ct.Count
	}

	// Obtener delitos por estado
	type CrimeStatusCount struct {
		Status string
		Count  int64
	}
	var crimesByStatus []CrimeStatusCount
	err = r.db.SelectContext(ctx, &crimesByStatus, `
		SELECT status, COUNT(*) as count
		FROM crimes
		WHERE deleted_at IS NULL
		GROUP BY status
	`)
	if err != nil {
		return nil, fmt.Errorf("error al obtener delitos por estado: %w", err)
	}

	// Convertir a map
	crimesByStatusMap := make(map[string]int64)
	for _, cs := range crimesByStatus {
		crimesByStatusMap[cs.Status] = cs.Count
	}

	// Obtener delitos por ubicación
	type CrimeLocationCount struct {
		Location string
		Count    int64
	}
	var crimesByLocation []CrimeLocationCount
	err = r.db.SelectContext(ctx, &crimesByLocation, `
		SELECT city as location, COUNT(*) as count
		FROM crimes
		WHERE deleted_at IS NULL
		GROUP BY city
	`)
	if err != nil {
		return nil, fmt.Errorf("error al obtener delitos por ubicación: %w", err)
	}

	// Convertir a map
	crimesByLocationMap := make(map[string]int64)
	for _, cl := range crimesByLocation {
		crimesByLocationMap[cl.Location] = cl.Count
	}

	// Obtener delitos por dirección
	type CrimeAddressCount struct {
		Address string
		Count   int64
	}
	var crimesByAddress []CrimeAddressCount
	err = r.db.SelectContext(ctx, &crimesByAddress, `
		SELECT CONCAT(address, ' ', COALESCE(address_number, '')) as address, COUNT(*) as count
		FROM crimes
		WHERE deleted_at IS NULL
		GROUP BY address, address_number
	`)
	if err != nil {
		return nil, fmt.Errorf("error al obtener delitos por dirección: %w", err)
	}

	// Convertir a map
	crimesByAddressMap := make(map[string]int64)
	for _, ca := range crimesByAddress {
		crimesByAddressMap[ca.Address] = ca.Count
	}

	return &entities.CrimeStats{
		TotalCrimes:      total,
		ActiveCrimes:     active,
		InactiveCrimes:   inactive,
		CrimesByType:     crimesByTypeMap,
		CrimesByStatus:   crimesByStatusMap,
		CrimesByLocation: crimesByLocationMap,
		CrimesByAddress:  crimesByAddressMap,
		LastUpdate:       time.Now(),
	}, nil
}
