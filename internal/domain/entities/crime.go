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

// Crime representa un delito
type Crime struct {
	ID          string     `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Type        string     `json:"type" db:"crime_type"`
	Status      string     `json:"status" db:"status"`
	Location    Location   `json:"location"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Location representa la ubicación geográfica de un delito
type Location struct {
	Latitude      float64 `json:"latitude" db:"latitude"`
	Longitude     float64 `json:"longitude" db:"longitude"`
	Address       string  `json:"address" db:"address"`
	AddressNumber *string `json:"address_number,omitempty" db:"address_number"`
	City          *string `json:"city,omitempty" db:"city"`
	Province      *string `json:"province,omitempty" db:"province"`
	Country       *string `json:"country,omitempty" db:"country"`
	ZipCode       *string `json:"zip_code,omitempty" db:"zip_code"`
}

// CrimeList representa una lista paginada de delitos
type CrimeList struct {
	Items []Crime `json:"items"`
	Total int64   `json:"total"`
}

// CrimeStats representa las estadísticas de delitos
type CrimeStats struct {
	TotalCrimes      int64            `json:"total_crimes"`
	ActiveCrimes     int64            `json:"active_crimes"`
	InactiveCrimes   int64            `json:"inactive_crimes"`
	CrimesByType     map[string]int64 `json:"crimes_by_type"`
	CrimesByStatus   map[string]int64 `json:"crimes_by_status"`
	CrimesByLocation map[string]int64 `json:"crimes_by_location"`
	CrimesByAddress  map[string]int64 `json:"crimes_by_address"`
	LastUpdate       time.Time        `json:"last_update"`
}
