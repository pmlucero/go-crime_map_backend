package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"go-crime_map_backend/internal/domain/entities"

	_ "github.com/lib/pq"
)

const (
	insertLocationQuery = `
		INSERT INTO locations (latitude, longitude, address, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`

	insertCrimeQuery = `
		INSERT INTO crimes (id, type, description, location_id, date, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
)

// PostgresCrimeRepository implementa la interfaz CrimeRepository usando PostgreSQL
type PostgresCrimeRepository struct {
	db *sql.DB
}

// NewPostgresCrimeRepository crea una nueva instancia del repositorio
func NewPostgresCrimeRepository(db *sql.DB) *PostgresCrimeRepository {
	return &PostgresCrimeRepository{
		db: db,
	}
}

// Create persiste un nuevo delito en la base de datos
func (r *PostgresCrimeRepository) Create(ctx context.Context, crime *entities.Crime) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error al iniciar la transacción: %w", err)
	}
	defer tx.Rollback()

	var locationID int64
	err = tx.QueryRowContext(ctx, insertLocationQuery,
		crime.Location.Latitude,
		crime.Location.Longitude,
		crime.Location.Address,
		crime.CreatedAt,
		crime.UpdatedAt,
	).Scan(&locationID)
	if err != nil {
		return fmt.Errorf("error al insertar la ubicación: %w", err)
	}

	_, err = tx.ExecContext(ctx, insertCrimeQuery,
		crime.ID,
		crime.Type,
		crime.Description,
		locationID,
		crime.Date,
		crime.CreatedAt,
		crime.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error al insertar el delito: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar la transacción: %w", err)
	}

	log.Printf("[PostgresCrimeRepository] Delito creado exitosamente - ID: %s, Tipo: %s", crime.ID, crime.Type)

	return nil
}

// GetByID obtiene un delito por su ID
func (r *PostgresCrimeRepository) GetByID(ctx context.Context, id string) (*entities.Crime, error) {
	var crime entities.Crime
	var location entities.Location
	var locationID int64

	err := r.db.QueryRowContext(ctx,
		`SELECT c.id, c.type, c.description, c.date, c.created_at, c.updated_at,
				l.id, l.latitude, l.longitude, l.address
		 FROM crimes c
		 JOIN locations l ON c.location_id = l.id
		 WHERE c.id = $1`,
		id,
	).Scan(
		&crime.ID,
		&crime.Type,
		&crime.Description,
		&crime.Date,
		&crime.CreatedAt,
		&crime.UpdatedAt,
		&locationID,
		&location.Latitude,
		&location.Longitude,
		&location.Address,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener el delito: %w", err)
	}

	crime.Location = location
	return &crime, nil
}

// GetAll obtiene todos los delitos
func (r *PostgresCrimeRepository) GetAll(ctx context.Context) ([]*entities.Crime, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT c.id, c.type, c.description, c.date, c.created_at, c.updated_at,
				l.id, l.latitude, l.longitude, l.address
		 FROM crimes c
		 JOIN locations l ON c.location_id = l.id
		 ORDER BY c.date DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("error al obtener los delitos: %w", err)
	}
	defer rows.Close()

	var crimes []*entities.Crime
	for rows.Next() {
		var crime entities.Crime
		var location entities.Location
		var locationID int64

		err := rows.Scan(
			&crime.ID,
			&crime.Type,
			&crime.Description,
			&crime.Date,
			&crime.CreatedAt,
			&crime.UpdatedAt,
			&locationID,
			&location.Latitude,
			&location.Longitude,
			&location.Address,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear el delito: %w", err)
		}

		crime.Location = location
		crimes = append(crimes, &crime)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar los delitos: %w", err)
	}

	return crimes, nil
}

// Update actualiza un delito existente
func (r *PostgresCrimeRepository) Update(ctx context.Context, crime *entities.Crime) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error al iniciar la transacción: %w", err)
	}
	defer tx.Rollback()

	// Obtener el ID de la ubicación actual
	var locationID int64
	err = tx.QueryRowContext(ctx,
		`SELECT location_id FROM crimes WHERE id = $1`,
		crime.ID,
	).Scan(&locationID)
	if err != nil {
		return fmt.Errorf("error al obtener el ID de la ubicación: %w", err)
	}

	// Actualizar la ubicación
	_, err = tx.ExecContext(ctx,
		`UPDATE locations 
		 SET latitude = $1, longitude = $2, address = $3
		 WHERE id = $4`,
		crime.Location.Latitude,
		crime.Location.Longitude,
		crime.Location.Address,
		locationID,
	)
	if err != nil {
		return fmt.Errorf("error al actualizar la ubicación: %w", err)
	}

	// Actualizar el delito
	_, err = tx.ExecContext(ctx,
		`UPDATE crimes 
		 SET type = $1, description = $2, date = $3
		 WHERE id = $4`,
		crime.Type,
		crime.Description,
		crime.Date,
		crime.ID,
	)
	if err != nil {
		return fmt.Errorf("error al actualizar el delito: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar la transacción: %w", err)
	}

	return nil
}

// Delete elimina un delito por su ID
func (r *PostgresCrimeRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error al iniciar la transacción: %w", err)
	}
	defer tx.Rollback()

	// Obtener el ID de la ubicación
	var locationID int64
	err = tx.QueryRowContext(ctx,
		`SELECT location_id FROM crimes WHERE id = $1`,
		id,
	).Scan(&locationID)
	if err != nil {
		return fmt.Errorf("error al obtener el ID de la ubicación: %w", err)
	}

	// Eliminar el delito
	_, err = tx.ExecContext(ctx,
		`DELETE FROM crimes WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("error al eliminar el delito: %w", err)
	}

	// Eliminar la ubicación
	_, err = tx.ExecContext(ctx,
		`DELETE FROM locations WHERE id = $1`,
		locationID,
	)
	if err != nil {
		return fmt.Errorf("error al eliminar la ubicación: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar la transacción: %w", err)
	}

	return nil
}

// DeleteAll elimina todos los registros de la base de datos
func (r *PostgresCrimeRepository) DeleteAll() error {
	query := `
		DELETE FROM crimes;
		DELETE FROM locations;
	`
	_, err := r.db.Exec(query)
	return err
}

// Close cierra la conexión a la base de datos
func (r *PostgresCrimeRepository) Close() error {
	return r.db.Close()
}
