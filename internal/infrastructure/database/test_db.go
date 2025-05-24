package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
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
		panic(fmt.Sprintf("Error al conectar a la base de datos de prueba: %v", err))
	}

	// Crear tablas
	if err := createTables(db); err != nil {
		panic(fmt.Sprintf("Error al crear tablas: %v", err))
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
		panic(fmt.Sprintf("Error al conectar a la base de datos de prueba: %v", err))
	}
	defer func() { _ = db.Close() }()

	// Eliminar registros
	if err := DeleteRecords(db); err != nil {
		panic(fmt.Sprintf("Error al eliminar registros: %v", err))
	}
}

// createTables crea las tablas necesarias para las pruebas
func createTables(db *sqlx.DB) error {
	/*
		queries := []string{
			// Eliminar trigger antes de eliminar la tabla
			`DROP TRIGGER IF EXISTS update_crimes_updated_at ON crimes`,
			// Eliminar índices antes de eliminar la tabla
			`DROP INDEX IF EXISTS idx_crimes_status`,
			`DROP INDEX IF EXISTS idx_crimes_type`,
			`DROP INDEX IF EXISTS idx_crimes_location`,
			// Eliminar la tabla
			`DROP TABLE IF EXISTS crimes CASCADE`,
			// Crear función para el trigger si no existe
			`DO $do$
			BEGIN
				IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'update_updated_at_column') THEN
					CREATE OR REPLACE FUNCTION update_updated_at_column()
					RETURNS TRIGGER AS $func$
					BEGIN
						NEW.updated_at = CURRENT_TIMESTAMP;
						RETURN NEW;
					END;
					$func$ language plpgsql;
				END IF;
			END;
			$do$`,
			// Crear la tabla
			`CREATE TABLE IF NOT EXISTS crimes (
				id SERIAL PRIMARY KEY,
				uuid VARCHAR(36) UNIQUE NOT NULL,
				title VARCHAR(255) NOT NULL,
				description TEXT NOT NULL,
				crime_type VARCHAR(50) NOT NULL,
				status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
				latitude DOUBLE PRECISION NOT NULL,
				longitude DOUBLE PRECISION NOT NULL,
				address TEXT NOT NULL,
				address_number VARCHAR(50),
				city VARCHAR(100),
				province VARCHAR(100),
				country VARCHAR(100),
				zip_code VARCHAR(20),
				created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP WITH TIME ZONE
			)`,
			// Crear los índices
			`CREATE INDEX IF NOT EXISTS idx_crimes_status ON crimes(status)`,
			`CREATE INDEX IF NOT EXISTS idx_crimes_type ON crimes(crime_type)`,
			`CREATE INDEX IF NOT EXISTS idx_crimes_location ON crimes USING GIST (ll_to_earth(latitude, longitude))`,
			// Crear el trigger
			`CREATE TRIGGER update_crimes_updated_at
				BEFORE UPDATE ON crimes
				FOR EACH ROW
				EXECUTE FUNCTION update_updated_at_column()`,
		}
	*/
	queries := []string{
		`TRUNCATE TABLE crimes RESTART IDENTITY CASCADE`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("error al ejecutar query %s: %w", query, err)
		}
	}

	return nil
}

// GetEnvOrDefault obtiene una variable de entorno o devuelve un valor por defecto
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// DeleteRecords elimina todos los registros de las tablas
func DeleteRecords(db interface{}) error {
	switch v := db.(type) {
	case *sqlx.DB:
		queries := []string{
			`TRUNCATE TABLE crimes RESTART IDENTITY CASCADE`,
		}
		for _, query := range queries {
			if _, err := v.Exec(query); err != nil {
				// Ignorar el error si la tabla no existe
				if !strings.Contains(err.Error(), "does not exist") {
					return fmt.Errorf("error al ejecutar query %s: %w", query, err)
				}
			}
		}
		return nil
	case *sql.DB:
		queries := []string{
			`TRUNCATE TABLE crimes RESTART IDENTITY CASCADE`,
		}
		for _, query := range queries {
			if _, err := v.Exec(query); err != nil {
				// Ignorar el error si la tabla no existe
				if !strings.Contains(err.Error(), "does not exist") {
					return fmt.Errorf("error al ejecutar query %s: %w", query, err)
				}
			}
		}
		return nil
	default:
		return fmt.Errorf("tipo de base de datos no soportado: %T", db)
	}
}
