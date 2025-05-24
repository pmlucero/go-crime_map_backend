package tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/interface/controllers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// coverage: 90%

// MockSecurityRepository es un mock para el repositorio de seguridad
type MockSecurityRepository struct {
	mock.Mock
}

func (m *MockSecurityRepository) ValidateAPIKey(ctx context.Context, apiKey string) (*entities.APIKey, error) {
	args := m.Called(ctx, apiKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.APIKey), args.Error(1)
}

func TestSecurityController_ValidateAPIKey_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockSecurityRepository)
	controller := controllers.NewSecurityController(mockRepo)

	validAPIKey := &entities.APIKey{
		Key:       "test-api-key",
		Status:    "active",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	mockRepo.On("ValidateAPIKey", mock.Anything, "test-api-key").Return(validAPIKey, nil)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/validate", nil)
	req.Header.Set("X-API-Key", "test-api-key")
	w := httptest.NewRecorder()

	controller.ValidateAPIKey(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestSecurityController_ValidateAPIKey_MissingKey(t *testing.T) {
	// Arrange
	mockRepo := new(MockSecurityRepository)
	controller := controllers.NewSecurityController(mockRepo)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/validate", nil)
	w := httptest.NewRecorder()

	controller.ValidateAPIKey(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "API key is required")
}

func TestSecurityController_ValidateAPIKey_InvalidKey(t *testing.T) {
	// Arrange
	mockRepo := new(MockSecurityRepository)
	controller := controllers.NewSecurityController(mockRepo)

	mockRepo.On("ValidateAPIKey", mock.Anything, "invalid-key").Return(nil, errors.New("invalid api key"))

	// Act
	req := httptest.NewRequest(http.MethodGet, "/validate", nil)
	req.Header.Set("X-API-Key", "invalid-key")
	w := httptest.NewRecorder()

	controller.ValidateAPIKey(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid API key")
	mockRepo.AssertExpectations(t)
}

func TestSecurityController_ValidateAPIKey_InactiveKey(t *testing.T) {
	// Arrange
	mockRepo := new(MockSecurityRepository)
	controller := controllers.NewSecurityController(mockRepo)

	inactiveKey := &entities.APIKey{
		Key:       "inactive-key",
		Status:    "inactive",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	mockRepo.On("ValidateAPIKey", mock.Anything, "inactive-key").Return(inactiveKey, nil)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/validate", nil)
	req.Header.Set("X-API-Key", "inactive-key")
	w := httptest.NewRecorder()

	controller.ValidateAPIKey(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "API key is inactive")
	mockRepo.AssertExpectations(t)
}

func TestSecurityController_ValidateAPIKey_ExpiredKey(t *testing.T) {
	// Arrange
	mockRepo := new(MockSecurityRepository)
	controller := controllers.NewSecurityController(mockRepo)

	expiredKey := &entities.APIKey{
		Key:       "expired-key",
		Status:    "active",
		ExpiresAt: time.Now().Add(-24 * time.Hour),
	}

	mockRepo.On("ValidateAPIKey", mock.Anything, "expired-key").Return(expiredKey, nil)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/validate", nil)
	req.Header.Set("X-API-Key", "expired-key")
	w := httptest.NewRecorder()

	controller.ValidateAPIKey(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "API key has expired")
	mockRepo.AssertExpectations(t)
}
