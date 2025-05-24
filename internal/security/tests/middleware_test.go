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
	now := time.Now()
	apiKey := &entities.APIKey{
		ID:           "test-user",
		Key:          "test-key",
		ExpiresAt:    now.Add(24 * time.Hour),
		CreatedAt:    now,
		UpdatedAt:    now,
		IsActiveBool: true,
	}

	_, err := db.Exec(`
		INSERT INTO api_keys (key, user_id, is_active, created_at, expires_at)
		VALUES (?, ?, ?, ?, ?)
	`, apiKey.Key, apiKey.ID, true, apiKey.CreatedAt, apiKey.ExpiresAt)
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
