package security_test

import (
	"database/sql"
	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/security"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// coverage: 85%

func setupMiddlewareTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("error abriendo base de datos: %v", err)
	}

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

func TestAPIKeyMiddleware_WithoutAPIKey(t *testing.T) {
	// Arrange
	db := setupMiddlewareTestDB(t)
	repo := security.NewRepository(db)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := security.APIKeyMiddleware(repo)(next)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperado status %d, obtenido %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestAPIKeyMiddleware_WithAPIKey(t *testing.T) {
	// Arrange
	db := setupMiddlewareTestDB(t)
	repo := security.NewRepository(db)

	// Crear una API key válida
	apiKey := &entities.APIKey{
		ID:        "test-id",
		Key:       "test-key",
		Status:    "active",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := db.Exec(`
		INSERT INTO api_keys (id, key, status, expires_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, apiKey.ID, apiKey.Key, apiKey.Status, apiKey.ExpiresAt, apiKey.CreatedAt, apiKey.UpdatedAt)
	if err != nil {
		t.Fatalf("error insertando API key: %v", err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := security.APIKeyMiddleware(repo)(next)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-Key", apiKey.Key)
	rec := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("esperado status %d, obtenido %d", http.StatusOK, rec.Code)
	}
}
