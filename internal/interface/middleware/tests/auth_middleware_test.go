package middleware_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/interface/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestAuthMiddleware_MissingAPIKey_ReturnsUnauthorized(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockSecurityRepository)

	router := gin.New()
	router.Use(middleware.AuthMiddleware(mockRepo))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "API Key requerida", response["error"])
	mockRepo.AssertExpectations(t)
}

func TestAuthMiddleware_ValidAPIKey_ReturnsOK(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockSecurityRepository)
	mockRepo.On("ValidateAPIKey", mock.Anything, "valid-key").Return(&entities.APIKey{
		ID:           "test-user",
		Key:          "valid-key",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		IsActiveBool: true,
	}, nil)

	router := gin.New()
	router.Use(middleware.AuthMiddleware(mockRepo))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "valid-key")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAuthMiddleware_ExpiredAPIKey_ReturnsUnauthorized(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockSecurityRepository)
	mockRepo.On("ValidateAPIKey", mock.Anything, "expired-key").Return(nil, errors.New("API Key expirada"))

	router := gin.New()
	router.Use(middleware.AuthMiddleware(mockRepo))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "expired-key")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "API Key expirada", response["error"])
	mockRepo.AssertExpectations(t)
}

func TestAuthMiddleware_InvalidAPIKey_ReturnsUnauthorized(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockSecurityRepository)
	mockRepo.On("ValidateAPIKey", mock.Anything, "invalid-key").Return(nil, errors.New("Error validando API Key"))

	router := gin.New()
	router.Use(middleware.AuthMiddleware(mockRepo))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "invalid-key")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Error validando API Key", response["error"])
	mockRepo.AssertExpectations(t)
}
