package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// SetupTestDB configura una base de datos de prueba
func SetupTestDB(t *testing.T) *sqlx.DB {
	// Obtener variables de entorno
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/crime_map_test?sslmode=disable"
	}

	// Conectar a la base de datos
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		t.Fatalf("Error al conectar a la base de datos de prueba: %v", err)
	}

	// Crear tablas
	if err := createTables(db); err != nil {
		t.Fatalf("Error al crear tablas: %v", err)
	}

	return db
}

// CleanupTestDB limpia la base de datos de prueba
func CleanupTestDB(t *testing.T) {
	// Obtener variables de entorno
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/crime_map_test?sslmode=disable"
	}

	// Conectar a la base de datos
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		t.Fatalf("Error al conectar a la base de datos de prueba: %v", err)
	}
	defer db.Close()

	// Eliminar registros
	if err := deleteRecords(db); err != nil {
		t.Fatalf("Error al eliminar registros: %v", err)
	}
}

// createTables crea las tablas necesarias para las pruebas
func createTables(db *sqlx.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS crimes (
			id VARCHAR(36) PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			crime_type VARCHAR(50) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			address TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_crimes_status ON crimes(status)`,
		`CREATE INDEX IF NOT EXISTS idx_crimes_type ON crimes(crime_type)`,
		`CREATE INDEX IF NOT EXISTS idx_crimes_location ON crimes USING GIST (ll_to_earth(latitude, longitude))`,
		`CREATE TRIGGER update_crimes_updated_at
			BEFORE UPDATE ON crimes
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column()`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("error al ejecutar query %s: %w", query, err)
		}
	}

	return nil
}

// deleteRecords elimina los registros de las tablas
func deleteRecords(db *sqlx.DB) error {
	queries := []string{
		`DELETE FROM crimes`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("error al ejecutar query %s: %w", query, err)
		}
	}

	return nil
}

// getEnvOrDefault obtiene una variable de entorno o devuelve un valor por defecto
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
