package controllers

import (
	"net/http"
	"time"

	"go-crime_map_backend/internal/domain/repositories"
)

// SecurityController maneja las operaciones relacionadas con la seguridad
type SecurityController struct {
	securityRepo repositories.SecurityRepository
}

// NewSecurityController crea una nueva instancia del controlador de seguridad
func NewSecurityController(securityRepo repositories.SecurityRepository) *SecurityController {
	return &SecurityController{
		securityRepo: securityRepo,
	}
}

// ValidateAPIKey valida una clave de API
func (c *SecurityController) ValidateAPIKey(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		http.Error(w, "API key is required", http.StatusUnauthorized)
		return
	}

	key, err := c.securityRepo.ValidateAPIKey(r.Context(), apiKey)
	if err != nil {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}

	if key == nil {
		http.Error(w, "API key not found", http.StatusUnauthorized)
		return
	}

	if key.Status != "active" {
		http.Error(w, "API key is inactive", http.StatusUnauthorized)
		return
	}

	if time.Now().After(key.ExpiresAt) {
		http.Error(w, "API key has expired", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
