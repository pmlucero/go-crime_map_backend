package security_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/security"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// coverage: 85%

func setupTestController(t *testing.T) (*security.Controller, *security.Repository, *sql.DB) {
	db := setupTestDB(t)
	// Crear tabla users
	_, err := db.Exec(`
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
	// Crear la tabla api_keys igual que en migraciones
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
	repo := security.NewRepository(db)
	controller := security.NewController(repo)
	return controller, repo, db
}

func TestController_GenerateAPIKey(t *testing.T) {
	// Arrange
	controller, _, _ := setupTestController(t)
	reqBody := map[string]string{
		"user_id": "test-user",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api-keys", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	// Act
	controller.GenerateAPIKeyHandler(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("esperado status %d, obtenido %d", http.StatusOK, rec.Code)
	}

	var response security.APIKeyResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("error decodificando respuesta: %v", err)
	}

	if response.Key == "" {
		t.Error("API Key generada está vacía")
	}
	if response.ID != "test-user" {
		t.Errorf("esperado user_id test-user, obtenido %s", response.ID)
	}
}

func TestController_RevokeAPIKey(t *testing.T) {
	// Arrange
	controller, repo, db := setupTestController(t)
	ctx := context.Background()

	// Primero creamos una API Key
	now := time.Now()
	apiKey := &entities.APIKey{
		ID:           "test-user",
		Key:          "test-key",
		ExpiresAt:    now.Add(24 * time.Hour),
		CreatedAt:    now,
		UpdatedAt:    now,
		IsActiveBool: true,
	}
	_, err := db.ExecContext(ctx, `
		INSERT INTO api_keys (key, user_id, is_active, created_at, expires_at)
		VALUES (?, ?, ?, ?, ?)
	`, apiKey.Key, apiKey.ID, apiKey.IsActiveBool, apiKey.CreatedAt, apiKey.ExpiresAt)
	if err != nil {
		t.Fatalf("error insertando API key: %v", err)
	}

	req := httptest.NewRequest("DELETE", "/api-keys?key=test-key", nil)
	rec := httptest.NewRecorder()

	// Act
	controller.RevokeAPIKeyHandler(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("esperado status %d, obtenido %d", http.StatusOK, rec.Code)
	}

	// Verificar que la API Key fue revocada
	storedKey, err := repo.GetAPIKey(ctx, "test-key")
	if err != nil {
		t.Fatalf("error obteniendo API Key: %v", err)
	}
	if storedKey.IsActiveBool != false {
		t.Error("API Key no fue revocada")
	}
}

func TestController_ListAPIKeys(t *testing.T) {
	// Arrange
	controller, _, db := setupTestController(t)
	ctx := context.Background()

	// Crear algunas API Keys para el usuario
	now := time.Now()
	apiKeys := []*entities.APIKey{
		{
			ID:           "test-user-1",
			Key:          "key1",
			ExpiresAt:    now.Add(24 * time.Hour),
			CreatedAt:    now,
			UpdatedAt:    now,
			IsActiveBool: true,
		},
		{
			ID:           "test-user-2",
			Key:          "key2",
			ExpiresAt:    now.Add(24 * time.Hour),
			CreatedAt:    now,
			UpdatedAt:    now,
			IsActiveBool: true,
		},
	}

	for _, key := range apiKeys {
		_, err := db.ExecContext(ctx, `
			INSERT INTO api_keys (key, user_id, is_active, created_at, expires_at)
			VALUES (?, ?, ?, ?, ?)
		`, key.Key, key.ID, key.IsActiveBool, key.CreatedAt, key.ExpiresAt)
		if err != nil {
			t.Fatalf("error insertando API key: %v", err)
		}
	}

	req := httptest.NewRequest("GET", "/api-keys?user_id=test-user", nil)
	rec := httptest.NewRecorder()

	// Act
	controller.ListAPIKeysHandler(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("esperado status %d, obtenido %d", http.StatusOK, rec.Code)
	}

	var response []entities.APIKey
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("error decodificando respuesta: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("esperado 2 API Keys, obtenido %d", len(response))
	}
}
