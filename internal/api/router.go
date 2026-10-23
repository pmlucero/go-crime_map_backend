package api

import (
	"go-crime_map_backend/internal/security"
	"net/http"
)

// RegisterRoutes registra todas las rutas de la API
func RegisterRoutes(mux *http.ServeMux, securityRepo *security.Repository) {
	// Crear el controlador de seguridad
	securityController := security.NewController(securityRepo)

	// Rutas de API Keys
	mux.HandleFunc("POST /api-keys", securityController.GenerateAPIKeyHandler)
	mux.HandleFunc("DELETE /api-keys", securityController.RevokeAPIKeyHandler)
	mux.HandleFunc("GET /api-keys", securityController.ListAPIKeysHandler)

	// TODO: Agregar el middleware de API Key a las rutas protegidas
	// Por ahora solo registramos las rutas de gestión de API Keys
}
