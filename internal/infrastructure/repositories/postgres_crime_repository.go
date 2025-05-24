package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go-crime_map_backend/internal/domain/entities"

	"github.com/google/uuid"
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
	crime.UUID = uuid.New().String()
	crime.CreatedAt = time.Now()
	crime.UpdatedAt = time.Now()
	crime.Status = string(entities.CrimeStatusActive)

	query := `
		INSERT INTO crimes (
			uuid, title, crime_type, description, latitude, longitude, status, 
			address, address_number, city, province, country, zip_code,
			created_at, updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, 
			$8, $9, $10, $11, $12, $13,
			$14, $15
		)
		RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		crime.UUID,
		crime.Title,
		crime.Type,
		crime.Description,
		crime.Location.Latitude,
		crime.Location.Longitude,
		crime.Status,
		crime.Location.Address,
		crime.Location.AddressNumber,
		crime.Location.City,
		crime.Location.Province,
		crime.Location.Country,
		crime.Location.ZipCode,
		crime.CreatedAt,
		crime.UpdatedAt,
	).Scan(&crime.ID)

	if err != nil {
		return fmt.Errorf("error al crear el delito: %w", err)
	}

	return nil
}

// List obtiene una lista paginada de delitos
func (r *PostgresCrimeRepository) List(ctx context.Context, page, limit int, startDate, endDate *time.Time, crimeType, status *string) ([]entities.Crime, int64, error) {
	query := `SELECT id, uuid, title, description, crime_type as type, status, latitude, longitude, address, address_number, city, province, country, zip_code, created_at, updated_at FROM crimes WHERE deleted_at IS NULL`
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
	defer func() { _ = rows.Close() }()

	var crimes []entities.Crime
	for rows.Next() {
		var crime entities.Crime
		var addressNumber, city, province, country, zipCode sql.NullString
		err := rows.Scan(
			&crime.ID,
			&crime.UUID,
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
			crime.Location.AddressNumber = addressNumber.String
		}
		if city.Valid {
			crime.Location.City = city.String
		}
		if province.Valid {
			crime.Location.Province = province.String
		}
		if country.Valid {
			crime.Location.Country = country.String
		}
		if zipCode.Valid {
			crime.Location.ZipCode = zipCode.String
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

// GetByID obtiene un delito por su ID numérico
func (r *PostgresCrimeRepository) GetByID(ctx context.Context, id int64) (*entities.Crime, error) {
	crime := &entities.Crime{}
	query := `
		SELECT id, uuid, type, description, latitude, longitude, status, created_at, updated_at, deleted_at
		FROM crimes
		WHERE id = $1`

	err := r.db.GetContext(ctx, crime, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("delito no encontrado con ID %d", id)
		}
		return nil, fmt.Errorf("error al obtener el delito: %w", err)
	}

	return crime, nil
}

// GetByUUID obtiene un delito por su UUID
func (r *PostgresCrimeRepository) GetByUUID(ctx context.Context, uuid string) (*entities.Crime, error) {
	crime := &entities.Crime{}
	query := `
		SELECT id, uuid, title, description, crime_type as type, status, 
			   latitude, longitude, address, address_number, city, province, country, zip_code,
			   created_at, updated_at, deleted_at
		FROM crimes
		WHERE uuid = $1`

	err := r.db.GetContext(ctx, crime, query, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("delito no encontrado con UUID %s", uuid)
		}
		return nil, fmt.Errorf("error al obtener el delito: %w", err)
	}

	return crime, nil
}

// Update actualiza un delito existente
func (r *PostgresCrimeRepository) Update(ctx context.Context, crime *entities.Crime) error {
	crime.UpdatedAt = time.Now()

	query := `
		UPDATE crimes
		SET title = $1, crime_type = $2, description = $3, latitude = $4, longitude = $5, 
		    status = $6, updated_at = $7
		WHERE id = $8 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(
		ctx,
		query,
		crime.Title,
		crime.Type,
		crime.Description,
		crime.Location.Latitude,
		crime.Location.Longitude,
		crime.Status,
		crime.UpdatedAt,
		crime.ID,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar el delito: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener las filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no se encontró el delito para actualizar")
	}

	return nil
}

// Delete elimina un delito por su ID
func (r *PostgresCrimeRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE crimes
		SET status = $1, deleted_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, "DELETED", id)
	if err != nil {
		return fmt.Errorf("error al eliminar el delito: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener las filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no se encontró el delito para eliminar")
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
	activeCrimes = crimesByStatus["ACTIVE"]
	inactiveCrimes = crimesByStatus["INACTIVE"]

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
