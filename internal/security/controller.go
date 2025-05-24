package security

import (
	"encoding/json"
	"net/http"
	"time"

	"go-crime_map_backend/internal/domain/entities"
)

// Controller maneja las operaciones HTTP relacionadas con API Keys
type Controller struct {
	repo *Repository
}

// NewController crea una nueva instancia del controlador
func NewController(repo *Repository) *Controller {
	return &Controller{repo: repo}
}

// GenerateAPIKeyRequest representa la solicitud para generar una nueva API Key
type GenerateAPIKeyRequest struct {
	UserID string `json:"user_id"`
}

// APIKeyResponse representa la respuesta con la API Key generada
type APIKeyResponse struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GenerateAPIKeyHandler maneja la generación de nuevas API Keys
func (c *Controller) GenerateAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	var req GenerateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decodificando la solicitud", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "user_id es requerido", http.StatusBadRequest)
		return
	}

	// Generar la API Key
	key, err := GenerateAPIKey()
	if err != nil {
		http.Error(w, "Error generando API Key", http.StatusInternalServerError)
		return
	}

	// Crear y almacenar la API Key
	now := time.Now()
	apiKey := &entities.APIKey{
		ID:        req.UserID,
		Key:       key,
		Status:    "active",
		ExpiresAt: now.Add(24 * time.Hour * 30), // 30 días
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := c.repo.CreateAPIKey(r.Context(), apiKey); err != nil {
		http.Error(w, "Error almacenando API Key", http.StatusInternalServerError)
		return
	}

	// Responder con la API Key generada
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(APIKeyResponse{
		ID:        apiKey.ID,
		Key:       apiKey.Key,
		Status:    apiKey.Status,
		ExpiresAt: apiKey.ExpiresAt,
		CreatedAt: apiKey.CreatedAt,
		UpdatedAt: apiKey.UpdatedAt,
	}); err != nil {
		http.Error(w, "Error codificando la respuesta", http.StatusInternalServerError)
	}
}

// RevokeAPIKeyHandler maneja la revocación de API Keys
func (c *Controller) RevokeAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key es requerido", http.StatusBadRequest)
		return
	}

	if err := c.repo.RevokeAPIKey(r.Context(), key); err != nil {
		http.Error(w, "Error revocando API Key", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ListAPIKeysHandler maneja la lista de API Keys
func (c *Controller) ListAPIKeysHandler(w http.ResponseWriter, r *http.Request) {
	apiKeys, err := c.repo.ListAPIKeys(r.Context())
	if err != nil {
		http.Error(w, "Error obteniendo API Keys", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(apiKeys); err != nil {
		http.Error(w, "Error codificando la respuesta", http.StatusInternalServerError)
	}
}
