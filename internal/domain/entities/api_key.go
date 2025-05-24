package entities

import "time"

// APIKey representa una clave de API
type APIKey struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Status    string    `json:"status"` // active, inactive
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsActive verifica si la clave de API está activa
func (a *APIKey) IsActive() bool {
	return a.Status == "active"
}

// IsExpired verifica si la clave de API ha expirado
func (a *APIKey) IsExpired() bool {
	return time.Now().After(a.ExpiresAt)
}
