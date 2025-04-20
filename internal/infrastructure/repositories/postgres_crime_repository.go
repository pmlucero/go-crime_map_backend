package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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
func (r *PostgresCrimeRepository) List(ctx context.Context, page, limit int, startDate, endDate *time.Time, crimeType, status *string) ([]entities.Crime, int64, error) {
	query := `SELECT id, title, description, crime_type as type, status, latitude, longitude, address, address_number, city, province, country, zip_code, created_at, updated_at FROM crimes WHERE deleted_at IS NULL`
	args := []interface{}{}

	// Agregar filtros
	if startDate != nil {
		query += " AND created_at >= $1"
		args = append(args, startDate)
	}
	if endDate != nil {
		query += " AND created_at <= $" + strconv.Itoa(len(args)+1)
		args = append(args, endDate)
	}
	if crimeType != nil {
		query += " AND crime_type = $" + strconv.Itoa(len(args)+1)
		args = append(args, crimeType)
	}
	if status != nil {
		query += " AND status = $" + strconv.Itoa(len(args)+1)
		args = append(args, status)
	}

	// Agregar paginación
	query += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, (page-1)*limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var crimes []entities.Crime
	for rows.Next() {
		var crime entities.Crime
		var addressNumber, city, province, country, zipCode sql.NullString
		err := rows.Scan(
			&crime.ID,
			&crime.Title,
			&crime.Description,
			&crime.Type,
			&crime.Status,
			&crime.Location.Latitude,
			&crime.Location.Longitude,
			&crime.Location.Address,
			&addressNumber,
			&city,
			&province,
			&country,
			&zipCode,
			&crime.CreatedAt,
			&crime.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if addressNumber.Valid {
			crime.Location.AddressNumber = &addressNumber.String
		}
		if city.Valid {
			crime.Location.City = &city.String
		}
		if province.Valid {
			crime.Location.Province = &province.String
		}
		if country.Valid {
			crime.Location.Country = &country.String
		}
		if zipCode.Valid {
			crime.Location.ZipCode = &zipCode.String
		}

		crimes = append(crimes, crime)
	}

	// Obtener el total de registros
	countQuery := `SELECT COUNT(*) FROM crimes WHERE deleted_at IS NULL`
	countArgs := []interface{}{}

	if startDate != nil {
		countQuery += " AND created_at >= $1"
		countArgs = append(countArgs, startDate)
	}
	if endDate != nil {
		countQuery += " AND created_at <= $" + strconv.Itoa(len(countArgs)+1)
		countArgs = append(countArgs, endDate)
	}
	if crimeType != nil {
		countQuery += " AND crime_type = $" + strconv.Itoa(len(countArgs)+1)
		countArgs = append(countArgs, crimeType)
	}
	if status != nil {
		countQuery += " AND status = $" + strconv.Itoa(len(countArgs)+1)
		countArgs = append(countArgs, status)
	}

	var total int64
	err = r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
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

	// Obtener delitos por tipo
	type crimeTypeCount struct {
		Type  string `db:"crime_type"`
		Count int64  `db:"count"`
	}
	var crimesByTypeRows []crimeTypeCount
	err = r.db.SelectContext(ctx, &crimesByTypeRows, `
		SELECT crime_type, COUNT(*) as count
		FROM crimes
		WHERE deleted_at IS NULL
		GROUP BY crime_type
	`)
	if err != nil {
		return nil, fmt.Errorf("error al obtener delitos por tipo: %w", err)
	}
	crimesByType := make(map[string]int64)
	for _, row := range crimesByTypeRows {
		crimesByType[row.Type] = row.Count
	}

	// Obtener delitos por estado
	type crimeStatusCount struct {
		Status string `db:"status"`
		Count  int64  `db:"count"`
	}
	var crimesByStatusRows []crimeStatusCount
	err = r.db.SelectContext(ctx, &crimesByStatusRows, `
		SELECT status, COUNT(*) as count
		FROM crimes
		WHERE deleted_at IS NULL
		GROUP BY status
	`)
	if err != nil {
		return nil, fmt.Errorf("error al obtener delitos por estado: %w", err)
	}
	crimesByStatus := make(map[string]int64)
	for _, row := range crimesByStatusRows {
		crimesByStatus[row.Status] = row.Count
	}

	// Obtener delitos por dirección
	type crimeAddressCount struct {
		Address string `db:"address"`
		Count   int64  `db:"count"`
	}
	var crimesByAddressRows []crimeAddressCount
	err = r.db.SelectContext(ctx, &crimesByAddressRows, `
		SELECT 
			CASE 
				WHEN address_number IS NOT NULL THEN address || ' ' || address_number
				ELSE address
			END as address,
			COUNT(*) as count
		FROM crimes
		WHERE deleted_at IS NULL
		GROUP BY 
			CASE 
				WHEN address_number IS NOT NULL THEN address || ' ' || address_number
				ELSE address
			END
	`)
	if err != nil {
		return nil, fmt.Errorf("error al obtener delitos por dirección: %w", err)
	}
	crimesByAddress := make(map[string]int64)
	for _, row := range crimesByAddressRows {
		crimesByAddress[row.Address] = row.Count
	}

	// Calcular delitos activos e inactivos
	var activeCrimes, inactiveCrimes int64
	for status, count := range crimesByStatus {
		if status == "ACTIVE" {
			activeCrimes = count
		} else if status == "INACTIVE" {
			inactiveCrimes = count
		}
	}

	return &entities.CrimeStats{
		TotalCrimes:     total,
		ActiveCrimes:    activeCrimes,
		InactiveCrimes:  inactiveCrimes,
		CrimesByType:    crimesByType,
		CrimesByStatus:  crimesByStatus,
		CrimesByAddress: crimesByAddress,
		LastUpdate:      time.Now(),
	}, nil
}
