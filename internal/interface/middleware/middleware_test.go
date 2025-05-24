package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
)

type mockSecurityRepo struct {
	validateAPIKeyFunc func(ctx context.Context, apiKey string) (*entities.APIKey, error)
}

func (m *mockSecurityRepo) ValidateAPIKey(ctx context.Context, apiKey string) (*entities.APIKey, error) {
	return m.validateAPIKeyFunc(ctx, apiKey)
}

func TestAPIKeyMiddleware_MissingAPIKey_ReturnsUnauthorized(t *testing.T) {
	// Arrange
	mockRepo := &mockSecurityRepo{
		validateAPIKeyFunc: func(ctx context.Context, apiKey string) (*entities.APIKey, error) {
			return nil, nil
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := APIKeyMiddleware(mockRepo)
	server := httptest.NewServer(middleware(handler))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Assert
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestAPIKeyMiddleware_ValidAPIKey_ReturnsOK(t *testing.T) {
	// Arrange
	validKey := &entities.APIKey{
		Key:          "valid-key",
		IsActiveBool: true,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	mockRepo := &mockSecurityRepo{
		validateAPIKeyFunc: func(ctx context.Context, apiKey string) (*entities.APIKey, error) {
			return validKey, nil
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := APIKeyMiddleware(mockRepo)
	server := httptest.NewServer(middleware(handler))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-API-Key", "valid-key")

	// Act
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestAPIKeyMiddleware_InactiveAPIKey_ReturnsUnauthorized(t *testing.T) {
	// Arrange
	inactiveKey := &entities.APIKey{
		Key:          "inactive-key",
		IsActiveBool: false,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	mockRepo := &mockSecurityRepo{
		validateAPIKeyFunc: func(ctx context.Context, apiKey string) (*entities.APIKey, error) {
			return inactiveKey, nil
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := APIKeyMiddleware(mockRepo)
	server := httptest.NewServer(middleware(handler))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-API-Key", "inactive-key")

	// Act
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Assert
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}
