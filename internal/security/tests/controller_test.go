package security_test

import (
	"bytes"
	"context"
	"encoding/json"
	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/security"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// coverage: 85%

func setupTestController(t *testing.T) (*security.Controller, *security.Repository) {
	db := setupTestDB(t)
	repo := security.NewRepository(db)
	controller := security.NewController(repo)
	return controller, repo
}

func TestController_GenerateAPIKey(t *testing.T) {
	// Arrange
	controller, _ := setupTestController(t)
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
	controller, repo := setupTestController(t)
	ctx := context.Background()

	// Primero creamos una API Key
	now := time.Now()
	apiKey := &entities.APIKey{
		ID:        "test-user",
		Key:       "test-key",
		Status:    "active",
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := repo.CreateAPIKey(ctx, apiKey); err != nil {
		t.Fatalf("error creando API Key: %v", err)
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
	if storedKey.Status != "inactive" {
		t.Error("API Key no fue revocada")
	}
}

func TestController_ListAPIKeys(t *testing.T) {
	// Arrange
	controller, repo := setupTestController(t)
	ctx := context.Background()

	// Crear algunas API Keys para el usuario
	now := time.Now()
	apiKeys := []*entities.APIKey{
		{
			ID:        "test-user-1",
			Key:       "key1",
			Status:    "active",
			ExpiresAt: now.Add(24 * time.Hour),
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "test-user-2",
			Key:       "key2",
			Status:    "active",
			ExpiresAt: now.Add(24 * time.Hour),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for _, key := range apiKeys {
		if err := repo.CreateAPIKey(ctx, key); err != nil {
			t.Fatalf("error creando API Key: %v", err)
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
