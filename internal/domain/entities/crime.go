package entities

import "time"

// CrimeStatus representa el estado de un delito
type CrimeStatus string

const (
	// CrimeStatusActive indica que el delito está activo
	CrimeStatusActive CrimeStatus = "ACTIVE"
	// CrimeStatusInactive indica que el delito está inactivo
	CrimeStatusInactive CrimeStatus = "INACTIVE"
	// CrimeStatusDeleted indica que el delito ha sido eliminado
	CrimeStatusDeleted CrimeStatus = "DELETED"
)

// Crime representa un delito reportado en el sistema
type Crime struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	Location    Location   `json:"location"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// Location representa la ubicación geográfica de un delito
type Location struct {
	Latitude  float64 `json:"latitude"`  // Latitud
	Longitude float64 `json:"longitude"` // Longitud
	Address   string  `json:"address"`   // Dirección descriptiva
}

// CrimeList representa una lista de delitos
type CrimeList struct {
	Crimes      []Crime `json:"crimes"`
	Total       int64   `json:"total"`
	Page        int     `json:"page"`
	Limit       int     `json:"limit"`
	TotalPages  int     `json:"total_pages"`
	HasNextPage bool    `json:"has_next_page"`
}

// CrimeStats representa las estadísticas de delitos
type CrimeStats struct {
	TotalCrimes      int64            `json:"total_crimes"`
	ActiveCrimes     int64            `json:"active_crimes"`
	InactiveCrimes   int64            `json:"inactive_crimes"`
	CrimesByType     map[string]int64 `json:"crimes_by_type"`
	CrimesByStatus   map[string]int64 `json:"crimes_by_status"`
	CrimesByLocation map[string]int64 `json:"crimes_by_location"`
}
