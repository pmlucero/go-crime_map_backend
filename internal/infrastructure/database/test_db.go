package database

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testDB *sql.DB
	once   sync.Once
)

// SetupTestDB configura la base de datos de test
func SetupTestDB(t *testing.T) *sql.DB {
	once.Do(func() {
		// Obtener variables de entorno para la base de datos de test
		dbHost := getEnvOrDefault("TEST_DB_HOST", "localhost")
		dbPort := getEnvOrDefault("TEST_DB_PORT", "5432")
		dbUser := getEnvOrDefault("TEST_DB_USER", "postgres")
		dbPassword := getEnvOrDefault("TEST_DB_PASSWORD", "postgres")
		dbName := getEnvOrDefault("TEST_DB_NAME", "crime_map_test")

		// Construir la cadena de conexión
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName)

		// Conectar a la base de datos
		var err error
		testDB, err = sql.Open("postgres", connStr)
		if err != nil {
			t.Fatalf("Error al conectar con la base de datos de test: %v", err)
		}

		// Verificar la conexión
		if err := testDB.Ping(); err != nil {
			t.Fatalf("Error al hacer ping a la base de datos de test: %v", err)
		}
	})

	// Limpiar la base de datos antes de crear las tablas
	CleanupTestDB(t)

	// Habilitar las extensiones necesarias
	_, err := testDB.Exec(`
		CREATE EXTENSION IF NOT EXISTS cube;
		CREATE EXTENSION IF NOT EXISTS earthdistance;
	`)
	if err != nil {
		t.Fatalf("Error al habilitar las extensiones: %v", err)
	}

	// Crear las tablas necesarias
	_, err = testDB.Exec(`
		CREATE TABLE IF NOT EXISTS locations (
			id SERIAL PRIMARY KEY,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			address TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS crimes (
			id UUID PRIMARY KEY,
			type VARCHAR(100) NOT NULL,
			description TEXT NOT NULL,
			location_id INTEGER NOT NULL REFERENCES locations(id),
			date TIMESTAMP WITH TIME ZONE NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_crimes_type ON crimes(type);
		CREATE INDEX IF NOT EXISTS idx_crimes_date ON crimes(date);
		CREATE INDEX IF NOT EXISTS idx_crimes_status ON crimes(status);
		CREATE INDEX IF NOT EXISTS idx_locations_coordinates ON locations USING GIST (ll_to_earth(latitude, longitude));

		-- Crear función para actualizar el timestamp updated_at
		CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ language 'plpgsql';

		-- Crear triggers para actualizar automáticamente el campo updated_at
		DROP TRIGGER IF EXISTS update_crimes_updated_at ON crimes;
		CREATE TRIGGER update_crimes_updated_at
			BEFORE UPDATE ON crimes
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();

		DROP TRIGGER IF EXISTS update_locations_updated_at ON locations;
		CREATE TRIGGER update_locations_updated_at
			BEFORE UPDATE ON locations
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();
	`)
	if err != nil {
		t.Fatalf("Error al crear las tablas: %v", err)
	}

	// Limpiar los datos de las tablas antes de cada test
	_, err = testDB.Exec(`
		TRUNCATE TABLE crimes CASCADE;
		TRUNCATE TABLE locations CASCADE;
	`)
	if err != nil {
		t.Fatalf("Error al limpiar los datos de las tablas: %v", err)
	}

	return testDB
}

// CleanupTestDB limpia la base de datos de test
func CleanupTestDB(t *testing.T) {
	if testDB == nil {
		return
	}

	// Eliminar todas las tablas y funciones
	_, err := testDB.Exec(`
		DROP TABLE IF EXISTS crimes CASCADE;
		DROP TABLE IF EXISTS locations CASCADE;
		DROP FUNCTION IF EXISTS update_updated_at_column CASCADE;
	`)
	if err != nil {
		t.Fatalf("Error al limpiar la base de datos de test: %v", err)
	}
}

// getEnvOrDefault obtiene una variable de entorno o devuelve un valor por defecto
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
