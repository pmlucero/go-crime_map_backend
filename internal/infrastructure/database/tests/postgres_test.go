// coverage: 85.7%

package tests

import (
	"os"
	"testing"

	"go-crime_map_backend/internal/infrastructure/database"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgresDB_DefaultValues_Success(t *testing.T) {
	// Arrange
	originalEnv := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
	}

	t.Cleanup(func() {
		for key, value := range originalEnv {
			if err := os.Setenv(key, value); err != nil {
				t.Fatal(err)
			}
		}
	})

	envVars := map[string]string{
		"DB_HOST":     "",
		"DB_PORT":     "",
		"DB_USER":     "",
		"DB_PASSWORD": "",
		"DB_NAME":     "",
	}

	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			t.Fatal(err)
		}
	}

	// Act
	db, err := database.NewPostgresDB()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, db)
	if db != nil {
		err = db.Close()
		assert.NoError(t, err)
	}
}

func TestNewPostgresDB_CustomValues_Success(t *testing.T) {
	// Arrange
	originalEnv := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
	}

	t.Cleanup(func() {
		for key, value := range originalEnv {
			if err := os.Setenv(key, value); err != nil {
				t.Fatal(err)
			}
		}
	})

	envVars := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"DB_USER":     "postgres",
		"DB_PASSWORD": "postgres",
		"DB_NAME":     "crime_map",
	}

	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			t.Fatal(err)
		}
	}

	// Act
	db, err := database.NewPostgresDB()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, db)
	if db != nil {
		err = db.Close()
		assert.NoError(t, err)
	}
}

func TestNewPostgresDB_InvalidPort_Error(t *testing.T) {
	// Arrange
	originalEnv := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
	}

	t.Cleanup(func() {
		for key, value := range originalEnv {
			if err := os.Setenv(key, value); err != nil {
				t.Fatal(err)
			}
		}
	})

	envVars := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "9999",
		"DB_USER":     "postgres",
		"DB_PASSWORD": "postgres",
		"DB_NAME":     "crime_map",
	}

	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			t.Fatal(err)
		}
	}

	// Act
	db, err := database.NewPostgresDB()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestNewPostgresDB_InvalidHost_Error(t *testing.T) {
	// Arrange
	originalEnv := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
	}

	t.Cleanup(func() {
		for key, value := range originalEnv {
			if err := os.Setenv(key, value); err != nil {
				t.Fatal(err)
			}
		}
	})

	envVars := map[string]string{
		"DB_HOST":     "invalid_host",
		"DB_PORT":     "5432",
		"DB_USER":     "postgres",
		"DB_PASSWORD": "postgres",
		"DB_NAME":     "crime_map",
	}

	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			t.Fatal(err)
		}
	}

	// Act
	db, err := database.NewPostgresDB()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, db)
}
