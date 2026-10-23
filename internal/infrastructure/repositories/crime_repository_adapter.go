package repositories

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// NewCrimeRepository crea una nueva instancia del repositorio de crímenes
func NewCrimeRepository(db *sql.DB) *PostgresCrimeRepository {
	// Convertir *sql.DB a *sqlx.DB
	dbx := sqlx.NewDb(db, "postgres")
	return NewPostgresCrimeRepository(dbx)
}
