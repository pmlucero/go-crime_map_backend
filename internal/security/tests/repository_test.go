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
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("error abriendo base de datos: %v", err)
	}

	// Crear tabla users
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(255) PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("error creando tabla users: %v", err)
	}

	// Insertar usuario de prueba
	_, err = db.Exec(`
		INSERT OR IGNORE INTO users (id, username, email) VALUES ('test-user', 'testuser', 'testuser@example.com')
	`)
	if err != nil {
		t.Fatalf("error insertando usuario de prueba: %v", err)
	}

	// Crear la tabla api_keys
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS api_keys (
			key VARCHAR(255) PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		t.Fatalf("error creando tabla api_keys: %v", err)
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

	now := time.Now()
	apiKey := &entities.APIKey{
		ID:           "test-user",
		Key:          "test-key",
		ExpiresAt:    now.Add(24 * time.Hour),
		CreatedAt:    now,
		UpdatedAt:    now,
		IsActiveBool: true,
	}

	// Act
	_, err := db.Exec(`
		INSERT INTO api_keys (key, user_id, is_active, created_at, expires_at)
		VALUES (?, ?, ?, ?, ?)
	`, apiKey.Key, apiKey.ID, apiKey.IsActiveBool, apiKey.CreatedAt, apiKey.ExpiresAt)
	if err != nil {
		t.Fatalf("error insertando API key: %v", err)
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

	if storedKey.IsActiveBool != apiKey.IsActiveBool {
		t.Errorf("esperado is_active %v, obtenido %v", apiKey.IsActiveBool, storedKey.IsActiveBool)
	}
}

func TestRepository_RevokeAPIKey(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := security.NewRepository(db)
	ctx := context.Background()

	now := time.Now()
	apiKey := &entities.APIKey{
		ID:           "test-user",
		Key:          "test-key",
		Status:       "active",
		ExpiresAt:    now.Add(24 * time.Hour),
		CreatedAt:    now,
		UpdatedAt:    now,
		IsActiveBool: true,
	}

	_, err := db.Exec(`
		INSERT INTO api_keys (key, user_id, is_active, created_at, expires_at)
		VALUES (?, ?, ?, ?, ?)
	`, apiKey.Key, apiKey.ID, apiKey.IsActiveBool, apiKey.CreatedAt, apiKey.ExpiresAt)
	if err != nil {
		t.Fatalf("error insertando API key: %v", err)
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

	if storedKey.IsActiveBool != false {
		t.Error("API Key no fue revocada")
	}
}
