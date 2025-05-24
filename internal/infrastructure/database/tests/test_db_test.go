package tests

import (
	"database/sql"
	"os"
	"testing"

	"go-crime_map_backend/internal/infrastructure/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupTestDB_ValidConfig_ReturnsConnection(t *testing.T) {
	// Arrange
	// No se necesitan preparaciones especiales

	// Act
	db := database.SetupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	// Assert
	assert.NotNil(t, db)
	err := db.Ping()
	assert.NoError(t, err)
}

func TestSetupTestDB_InvalidConfig_PanicsWithError(t *testing.T) {
	// Arrange
	originalURL := os.Getenv("TEST_DATABASE_URL")
	defer func() {
		if err := os.Setenv("TEST_DATABASE_URL", originalURL); err != nil {
			t.Fatal(err)
		}
	}()
	if err := os.Setenv("TEST_DATABASE_URL", "postgres://invalid:invalid@localhost:1111/invalid"); err != nil {
		t.Fatal(err)
	}

	// Act & Assert
	defer func() {
		if r := recover(); r != nil {
			panicMsg, ok := r.(string)
			assert.True(t, ok)
			assert.Contains(t, panicMsg, "Error al conectar a la base de datos de prueba")
		} else {
			t.Error("Se esperaba un panic")
		}
	}()

	database.SetupTestDB(t)
}

func TestCleanupTestDB_WithData_CleansAllTables(t *testing.T) {
	// Arrange
	db := database.SetupTestDB(t)
	require.NotNil(t, db)
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	_, err := db.Exec(`INSERT INTO crimes (uuid, title, description, crime_type, status, latitude, longitude, address) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		"550e8400-e29b-41d4-a716-446655440000", // UUID de ejemplo
		"Test Crime", "Test Description", "ROBO", "ACTIVE", 40.7128, -74.0060, "Test Address")
	require.NoError(t, err)

	// Act
	database.CleanupTestDB(t)

	// Assert
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM crimes").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestCleanupTestDB_InvalidConnection_PanicsWithError(t *testing.T) {
	// Arrange
	originalURL := os.Getenv("TEST_DATABASE_URL")
	defer func() {
		if err := os.Setenv("TEST_DATABASE_URL", originalURL); err != nil {
			t.Fatal(err)
		}
	}()
	if err := os.Setenv("TEST_DATABASE_URL", "postgres://invalid:invalid@localhost:1111/invalid"); err != nil {
		t.Fatal(err)
	}

	// Act & Assert
	defer func() {
		if r := recover(); r != nil {
			panicMsg, ok := r.(string)
			assert.True(t, ok)
			assert.Contains(t, panicMsg, "Error al conectar a la base de datos de prueba")
		} else {
			t.Error("Se esperaba un panic")
		}
	}()

	database.CleanupTestDB(t)
}

func TestDeleteRecords_WithSqlxDB_DeletesAllRecords(t *testing.T) {
	// Arrange
	db := database.SetupTestDB(t)
	require.NotNil(t, db)
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	// Limpiar la base de datos antes de comenzar
	err := database.DeleteRecords(db)
	require.NoError(t, err)

	// Insertar usando la secuencia automática y un UUID
	_, err = db.Exec(`INSERT INTO crimes (uuid, title, description, crime_type, status, latitude, longitude, address) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		"550e8400-e29b-41d4-a716-446655440000", // UUID de ejemplo
		"Test Crime", "Test Description", "ROBO", "ACTIVE", 40.7128, -74.0060, "Test Address")
	require.NoError(t, err)

	// Act
	err = database.DeleteRecords(db)

	// Assert
	assert.NoError(t, err)

	// Verificar que no hay registros
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM crimes").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	// Verificar que podemos insertar un nuevo registro sin conflictos
	_, err = db.Exec(`INSERT INTO crimes (uuid, title, description, crime_type, status, latitude, longitude, address) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		"550e8400-e29b-41d4-a716-446655440001", // UUID diferente
		"New Crime", "New Description", "ROBO", "ACTIVE", 40.7128, -74.0060, "New Address")
	assert.NoError(t, err, "Debería poder insertar un nuevo registro sin conflictos")
}

func TestDeleteRecords_WithSqlDB_DeletesAllRecords(t *testing.T) {
	// Arrange
	sqlDB, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/crime_map_test?sslmode=disable")
	require.NoError(t, err)
	defer func() { _ = sqlDB.Close() }()

	// Act
	err = database.DeleteRecords(sqlDB)

	// Assert
	assert.NoError(t, err)
}

func TestDeleteRecords_InvalidDBType_ReturnsError(t *testing.T) {
	// Arrange
	invalidDB := "invalid"

	// Act
	err := database.DeleteRecords(invalidDB)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tipo de base de datos no soportado")
}

func TestGetEnvOrDefault_ExistingVar_ReturnsValue(t *testing.T) {
	// Arrange
	key := "TEST_ENV_VAR"
	originalValue := os.Getenv(key)
	defer func() {
		if err := os.Setenv(key, originalValue); err != nil {
			t.Fatal(err)
		}
	}()
	if err := os.Setenv(key, "test_value"); err != nil {
		t.Fatal(err)
	}

	// Act
	value := database.GetEnvOrDefault(key, "default_value")

	// Assert
	assert.Equal(t, "test_value", value)
}

func TestGetEnvOrDefault_NonExistentVar_ReturnsDefault(t *testing.T) {
	// Arrange
	key := "NON_EXISTENT_VAR"

	// Act
	value := database.GetEnvOrDefault(key, "default_value")

	// Assert
	assert.Equal(t, "default_value", value)
}

// TODO: Revisar si es necesario
// Se comenta porque no se crea la tabla crimes por codigo
// sino por el script de la base de datos
/*
func TestCreateTables_NewDB_CreatesTablesAndIndexes(t *testing.T) {
	// Arrange
	db := database.SetupTestDB(t)
	defer database.CleanupTestDB(t)

	// Act & Assert
	// Verificar que la tabla crimes existe
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = 'crimes'
		)
	`).Scan(&exists)

	assert.NoError(t, err)
	assert.True(t, exists)

	// Verificar que los índices existen
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM pg_indexes
			WHERE indexname = 'idx_crimes_status'
		)
	`).Scan(&exists)

	assert.NoError(t, err)
	assert.True(t, exists)

	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM pg_indexes
			WHERE indexname = 'idx_crimes_type'
		)
	`).Scan(&exists)

	assert.NoError(t, err)
	assert.True(t, exists)

	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM pg_indexes
			WHERE indexname = 'idx_crimes_location'
		)
	`).Scan(&exists)

	assert.NoError(t, err)
	assert.True(t, exists)
}*/
