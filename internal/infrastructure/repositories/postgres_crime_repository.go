package repositories

import (
	"context"
	"database/sql"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
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
	// Convertir fechas a UTC
	now := time.Now().UTC()
	crime.CreatedAt = now
	crime.UpdatedAt = now
	crime.Date = crime.Date.UTC()

	// Primero insertar la ubicación
	var locationID int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO locations (latitude, longitude, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, crime.Location.Latitude, crime.Location.Longitude, crime.Location.Address, now, now).Scan(&locationID)
	if err != nil {
		return err
	}

	// Luego insertar el delito
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO crimes (id, type, description, location_id, date, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, crime.ID, crime.Type, crime.Description, locationID, crime.Date, crime.Status, now, now)
	return err
}

// GetByID obtiene un delito por su ID
func (r *PostgresCrimeRepository) GetByID(ctx context.Context, id string) (*entities.Crime, error) {
	var crime entities.Crime
	var locationID int
	err := r.db.QueryRowContext(ctx, `
		SELECT c.id, c.type, c.description, c.date, c.status, c.created_at, c.updated_at,
			   l.id, l.latitude, l.longitude, l.address
		FROM crimes c
		JOIN locations l ON c.location_id = l.id
		WHERE c.id = $1
	`, id).Scan(
		&crime.ID, &crime.Type, &crime.Description, &crime.Date, &crime.Status,
		&crime.CreatedAt, &crime.UpdatedAt,
		&locationID, &crime.Location.Latitude, &crime.Location.Longitude, &crime.Location.Address,
	)
	if err == sql.ErrNoRows {
		return nil, repositories.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Convertir fechas a UTC
	crime.Date = crime.Date.UTC()
	crime.CreatedAt = crime.CreatedAt.UTC()
	crime.UpdatedAt = crime.UpdatedAt.UTC()

	return &crime, nil
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
		var locationID int
		err := rows.Scan(
			&crime.ID, &crime.Type, &crime.Description, &crime.Date, &crime.Status,
			&crime.CreatedAt, &crime.UpdatedAt,
			&locationID, &crime.Location.Latitude, &crime.Location.Longitude, &crime.Location.Address,
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
	// Primero actualizar la ubicación
	_, err := r.db.ExecContext(ctx, `
		UPDATE locations l
		SET latitude = $1, longitude = $2, address = $3
		FROM crimes c
		WHERE c.location_id = l.id AND c.id = $4
	`, crime.Location.Latitude, crime.Location.Longitude, crime.Location.Address, crime.ID)
	if err != nil {
		return err
	}

	// Luego actualizar el delito
	_, err = r.db.ExecContext(ctx, `
		UPDATE crimes
		SET type = $1, description = $2, date = $3, status = $4
		WHERE id = $5
	`, crime.Type, crime.Description, crime.Date, crime.Status, crime.ID)
	return err
}

// Delete elimina un delito por su ID
func (r *PostgresCrimeRepository) Delete(ctx context.Context, id string) error {
	// La eliminación en cascada se maneja en la base de datos
	_, err := r.db.ExecContext(ctx, "DELETE FROM crimes WHERE id = $1", id)
	return err
}

// List obtiene una lista de delitos con los filtros especificados
func (r *PostgresCrimeRepository) List(ctx context.Context, filter repositories.ListCrimesFilter) ([]*entities.Crime, error) {
	query := `
		SELECT c.id, c.type, c.description, c.date, c.status, c.created_at, c.updated_at,
			   l.id, l.latitude, l.longitude, l.address
		FROM crimes c
		JOIN locations l ON c.location_id = l.id
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	if filter.Type != "" {
		query += ` AND c.type = $` + string(rune('0'+argCount))
		args = append(args, filter.Type)
		argCount++
	}

	if filter.Status != "" {
		query += ` AND c.status = $` + string(rune('0'+argCount))
		args = append(args, filter.Status)
		argCount++
	}

	if !filter.StartDate.IsZero() {
		query += ` AND c.date >= $` + string(rune('0'+argCount))
		args = append(args, filter.StartDate)
		argCount++
	}

	if !filter.EndDate.IsZero() {
		query += ` AND c.date <= $` + string(rune('0'+argCount))
		args = append(args, filter.EndDate)
		argCount++
	}

	if filter.Latitude != 0 && filter.Longitude != 0 && filter.Radius > 0 {
		query += ` AND earth_distance(ll_to_earth(l.latitude, l.longitude), ll_to_earth($` + string(rune('0'+argCount)) + `, $` + string(rune('0'+argCount+1)) + `)) <= $` + string(rune('0'+argCount+2))
		args = append(args, filter.Latitude, filter.Longitude, filter.Radius*1000) // Convertir km a metros
		argCount += 3
	}

	query += ` ORDER BY c.created_at DESC LIMIT $` + string(rune('0'+argCount)) + ` OFFSET $` + string(rune('0'+argCount+1))
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var crimes []*entities.Crime
	for rows.Next() {
		var crime entities.Crime
		var locationID int
		err := rows.Scan(
			&crime.ID, &crime.Type, &crime.Description, &crime.Date, &crime.Status,
			&crime.CreatedAt, &crime.UpdatedAt,
			&locationID, &crime.Location.Latitude, &crime.Location.Longitude, &crime.Location.Address,
		)
		if err != nil {
			return nil, err
		}
		crimes = append(crimes, &crime)
	}
	return crimes, nil
}
