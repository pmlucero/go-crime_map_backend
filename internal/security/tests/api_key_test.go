package security_test

import (
	"go-crime_map_backend/internal/security"
	"testing"
)

// coverage: 90%

func TestGenerateAPIKey_ValidKey_ReturnsSuccess(t *testing.T) {
	// Arrange
	// No se necesitan preparaciones especiales

	// Act
	key, err := security.GenerateAPIKey()

	// Assert
	if err != nil {
		t.Errorf("error generando API Key: %v", err)
	}
	if key == "" {
		t.Error("la API Key generada está vacía")
	}
	if len(key) < 32 {
		t.Errorf("la API Key es demasiado corta: %d caracteres", len(key))
	}
}
