package tests

import (
	"database/sql"
	"os"
	"testing"

	"go-crime_map_backend/internal/infrastructure/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupTestDB(t *testing.T) {
	t.Run("Configuración exitosa", func(t *testing.T) {
		db := database.SetupTestDB(t)
		assert.NotNil(t, db)
		defer db.Close()

		// Verificar que la conexión funciona
		err := db.Ping()
		assert.NoError(t, err)
	})

	t.Run("Error de conexión", func(t *testing.T) {
		// Guardar la URL original
		originalURL := os.Getenv("TEST_DATABASE_URL")
		defer os.Setenv("TEST_DATABASE_URL", originalURL)

		// Establecer una URL inválida
		os.Setenv("TEST_DATABASE_URL", "postgres://invalid:invalid@localhost:1111/invalid")

		// La función debería hacer panic, así que usamos recover
		defer func() {
			if r := recover(); r != nil {
				// Verificar que el mensaje de error es el esperado
				panicMsg, ok := r.(string)
				assert.True(t, ok)
				assert.Contains(t, panicMsg, "Error al conectar a la base de datos de prueba")
			} else {
				t.Error("Se esperaba un panic")
			}
		}()

		database.SetupTestDB(t)
	})
}

func TestCleanupTestDB(t *testing.T) {
	t.Run("Limpieza exitosa", func(t *testing.T) {
		db := database.SetupTestDB(t)
		require.NotNil(t, db)
		defer db.Close()

		// Insertar algunos datos
		_, err := db.Exec(`INSERT INTO crimes (id, title, description, crime_type, status, latitude, longitude, address) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			"test-id", "Test Crime", "Test Description", "ROBO", "ACTIVE", 40.7128, -74.0060, "Test Address")
		require.NoError(t, err)

		// Ejecutar la limpieza
		database.CleanupTestDB(t)

		// Verificar que los datos fueron eliminados
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM crimes").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("Error de conexión", func(t *testing.T) {
		// Guardar la URL original
		originalURL := os.Getenv("TEST_DATABASE_URL")
		defer os.Setenv("TEST_DATABASE_URL", originalURL)

		// Establecer una URL inválida
		os.Setenv("TEST_DATABASE_URL", "postgres://invalid:invalid@localhost:1111/invalid")

		// La función debería hacer panic, así que usamos recover
		defer func() {
			if r := recover(); r != nil {
				// Verificar que el mensaje de error es el esperado
				panicMsg, ok := r.(string)
				assert.True(t, ok)
				assert.Contains(t, panicMsg, "Error al conectar a la base de datos de prueba")
			} else {
				t.Error("Se esperaba un panic")
			}
		}()

		database.CleanupTestDB(t)
	})
}

func TestDeleteRecords(t *testing.T) {
	t.Run("Eliminar registros con *sqlx.DB", func(t *testing.T) {
		db := database.SetupTestDB(t)
		require.NotNil(t, db)
		defer db.Close()

		// Insertar algunos datos
		_, err := db.Exec(`INSERT INTO crimes (id, title, description, crime_type, status, latitude, longitude, address) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			"test-id", "Test Crime", "Test Description", "ROBO", "ACTIVE", 40.7128, -74.0060, "Test Address")
		require.NoError(t, err)

		// Eliminar registros
		err = database.DeleteRecords(db)
		assert.NoError(t, err)

		// Verificar que los datos fueron eliminados
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM crimes").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("Eliminar registros con *sql.DB", func(t *testing.T) {
		sqlDB, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/crime_map_test?sslmode=disable")
		require.NoError(t, err)
		defer sqlDB.Close()

		// Eliminar registros
		err = database.DeleteRecords(sqlDB)
		assert.NoError(t, err)
	})

	t.Run("Error con tipo no soportado", func(t *testing.T) {
		err := database.DeleteRecords("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tipo de base de datos no soportado")
	})
}

func TestGetEnvOrDefault(t *testing.T) {
	t.Run("Obtener valor existente", func(t *testing.T) {
		// Guardar el valor original
		key := "TEST_ENV_VAR"
		originalValue := os.Getenv(key)
		defer os.Setenv(key, originalValue)

		// Establecer un valor de prueba
		os.Setenv(key, "test_value")

		value := database.GetEnvOrDefault(key, "default_value")
		assert.Equal(t, "test_value", value)
	})

	t.Run("Obtener valor por defecto", func(t *testing.T) {
		value := database.GetEnvOrDefault("NON_EXISTENT_VAR", "default_value")
		assert.Equal(t, "default_value", value)
	})
}

func TestCreateTables(t *testing.T) {
	t.Run("Crear tablas exitosamente", func(t *testing.T) {
		db := database.SetupTestDB(t)
		defer database.CleanupTestDB(t)

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
	})
}
