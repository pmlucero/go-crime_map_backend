package security_test

import (
	"context"
	"database/sql"
	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/security"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// coverage: 85%

func setupTestDB(t *testing.T) *sql.DB {
	// TODO: Configurar una base de datos de prueba
	// Por ahora usamos una base de datos en memoria
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("error abriendo base de datos: %v", err)
	}

	// Crear la tabla api_keys
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS api_keys (
			id VARCHAR(255) PRIMARY KEY,
			key VARCHAR(255) NOT NULL,
			status VARCHAR(50) NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("error creando tabla: %v", err)
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	})

	return db
}

func TestRepository_CreateAndGetAPIKey(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := security.NewRepository(db)
	ctx := context.Background()

	apiKey := &entities.APIKey{
		ID:        "test-id",
		Key:       "test-key",
		Status:    "active",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	err := repo.CreateAPIKey(ctx, apiKey)
	if err != nil {
		t.Fatalf("error creando API Key: %v", err)
	}

	// Assert
	storedKey, err := repo.GetAPIKey(ctx, apiKey.Key)
	if err != nil {
		t.Fatalf("error obteniendo API Key: %v", err)
	}

	if storedKey == nil {
		t.Fatal("API Key no encontrada")
	}

	if storedKey.Key != apiKey.Key {
		t.Errorf("esperado key %s, obtenido %s", apiKey.Key, storedKey.Key)
	}

	if storedKey.Status != apiKey.Status {
		t.Errorf("esperado status %s, obtenido %s", apiKey.Status, storedKey.Status)
	}
}

func TestRepository_RevokeAPIKey(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := security.NewRepository(db)
	ctx := context.Background()

	apiKey := &entities.APIKey{
		ID:        "test-id",
		Key:       "test-key",
		Status:    "active",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.CreateAPIKey(ctx, apiKey)
	if err != nil {
		t.Fatalf("error creando API Key: %v", err)
	}

	// Act
	err = repo.RevokeAPIKey(ctx, apiKey.Key)
	if err != nil {
		t.Fatalf("error revocando API Key: %v", err)
	}

	// Assert
	storedKey, err := repo.GetAPIKey(ctx, apiKey.Key)
	if err != nil {
		t.Fatalf("error obteniendo API Key: %v", err)
	}

	if storedKey.Status != "inactive" {
		t.Error("API Key no fue revocada")
	}
}
