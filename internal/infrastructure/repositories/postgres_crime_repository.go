package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-crime_map_backend/internal/domain/entities"
)

// PostgresCrimeRepository implementa el repositorio de delitos usando PostgreSQL
type PostgresCrimeRepository struct {
	db *sql.DB
}

// NewPostgresCrimeRepository crea una nueva instancia del repositorio
func NewPostgresCrimeRepository(db *sql.DB) *PostgresCrimeRepository {
	return &PostgresCrimeRepository{
		db: db,
	}
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

// GetAll obtiene todos los delitos
func (r *PostgresCrimeRepository) GetAll(ctx context.Context) ([]*entities.Crime, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id, c.type, c.description, c.date, c.status, c.created_at, c.updated_at,
			   l.id, l.latitude, l.longitude, l.address
		FROM crimes c
		JOIN locations l ON c.location_id = l.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var crimes []*entities.Crime
	for rows.Next() {
		var crime entities.Crime
		err := rows.Scan(
			&crime.ID, &crime.Type, &crime.Description, &crime.Status,
			&crime.CreatedAt, &crime.UpdatedAt,
			&crime.Location.Latitude, &crime.Location.Longitude, &crime.Location.Address,
		)
		if err != nil {
			return nil, err
		}
		crimes = append(crimes, &crime)
	}
	return crimes, nil
}

// Update actualiza un delito existente
func (r *PostgresCrimeRepository) Update(ctx context.Context, crime *entities.Crime) error {
	query := `
		UPDATE crimes
		SET title = $1, description = $2, crime_type = $3, status = $4,
			latitude = $5, longitude = $6, address = $7, updated_at = $8
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
		crime.UpdatedAt,
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
		return fmt.Errorf("delito no encontrado o ya eliminado")
	}

	return nil
}

// Delete elimina un delito por su ID
func (r *PostgresCrimeRepository) Delete(ctx context.Context, id string) error {
	now := time.Now()
	query := `
		UPDATE crimes
		SET status = 'deleted', deleted_at = $1, updated_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, now, id)
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
func (r *PostgresCrimeRepository) List(ctx context.Context, page, limit int, startDate, endDate *time.Time, crimeType, status *string) ([]entities.Crime, int64, error) {
	offset := (page - 1) * limit

	// Construir la consulta base
	query := `
		SELECT c.id, c.title, c.description, c.crime_type, c.status, 
		       c.created_at, c.updated_at, c.deleted_at,
		       c.latitude, c.longitude, c.address
		FROM crimes c
		WHERE c.deleted_at IS NULL
	`
	countQuery := `SELECT COUNT(*) FROM crimes WHERE deleted_at IS NULL`

	// Agregar condiciones de filtro
	var conditions []string
	var args []interface{}
	argCount := 1

	if startDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argCount))
		args = append(args, startDate)
		argCount++
	}

	if endDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argCount))
		args = append(args, endDate)
		argCount++
	}

	if crimeType != nil {
		conditions = append(conditions, fmt.Sprintf("crime_type = $%d", argCount))
		args = append(args, *crimeType)
		argCount++
	}

	if status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *status)
		argCount++
	}

	// Agregar condiciones a las consultas
	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		query += whereClause
		countQuery += whereClause
	}

	// Agregar ordenamiento y paginación
	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argCount) + ` OFFSET $` + strconv.Itoa(argCount+1)
	args = append(args, limit, offset)

	// Ejecutar consulta de conteo
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error al contar delitos: %w", err)
	}

	// Ejecutar consulta principal
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error al listar delitos: %w", err)
	}
	defer rows.Close()

	var crimes []entities.Crime
	for rows.Next() {
		var crime entities.Crime
		err := rows.Scan(
			&crime.ID, &crime.Title, &crime.Description, &crime.Type, &crime.Status,
			&crime.CreatedAt, &crime.UpdatedAt,
			&crime.DeletedAt,
			&crime.Location.Latitude, &crime.Location.Longitude, &crime.Location.Address,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error al escanear delito: %w", err)
		}
		crimes = append(crimes, crime)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error al iterar delitos: %w", err)
	}

	return crimes, total, nil
}

// GetStats obtiene estadísticas sobre los delitos
func (r *PostgresCrimeRepository) GetStats(ctx context.Context) (*entities.CrimeStats, error) {
	// Obtener estadísticas básicas
	var stats entities.CrimeStats
	err := r.db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as total_crimes,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_crimes,
			COUNT(CASE WHEN status = 'inactive' THEN 1 END) as inactive_crimes
		FROM crimes
		WHERE deleted_at IS NULL
	`).Scan(&stats.TotalCrimes, &stats.ActiveCrimes, &stats.InactiveCrimes)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estadísticas básicas: %w", err)
	}

	// Obtener delitos por tipo
	rows, err := r.db.QueryContext(ctx, `
		SELECT crime_type, COUNT(*) as count
		FROM crimes
		WHERE deleted_at IS NULL
		GROUP BY crime_type
	`)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estadísticas por tipo: %w", err)
	}
	defer rows.Close()

	stats.CrimesByType = make(map[string]int64)
	for rows.Next() {
		var crimeType string
		var count int64
		if err := rows.Scan(&crimeType, &count); err != nil {
			return nil, fmt.Errorf("error al escanear estadísticas por tipo: %w", err)
		}
		stats.CrimesByType[crimeType] = count
	}

	// Obtener delitos por estado
	rows, err = r.db.QueryContext(ctx, `
		SELECT status, COUNT(*) as count
		FROM crimes
		WHERE deleted_at IS NULL
		GROUP BY status
	`)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estadísticas por estado: %w", err)
	}
	defer rows.Close()

	stats.CrimesByStatus = make(map[string]int64)
	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("error al escanear estadísticas por estado: %w", err)
		}
		stats.CrimesByStatus[status] = count
	}

	// Obtener delitos por ubicación
	rows, err = r.db.QueryContext(ctx, `
		SELECT address, COUNT(*) as count
		FROM crimes
		WHERE deleted_at IS NULL
		GROUP BY address
	`)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estadísticas por ubicación: %w", err)
	}
	defer rows.Close()

	stats.CrimesByLocation = make(map[string]int64)
	for rows.Next() {
		var address string
		var count int64
		if err := rows.Scan(&address, &count); err != nil {
			return nil, fmt.Errorf("error al escanear estadísticas por ubicación: %w", err)
		}
		stats.CrimesByLocation[address] = count
	}

	return &stats, nil
}
