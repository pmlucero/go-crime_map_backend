package entities

import (
	"errors"
	"time"
)

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

// CrimeType representa el tipo de delito
type CrimeType string

const (
	// CrimeTypeRobo representa un robo
	CrimeTypeRobo CrimeType = "ROBO"
	// CrimeTypeAgresion representa una agresión
	CrimeTypeAgresion CrimeType = "AGRESION"
	// CrimeTypeHurto representa un hurto
	CrimeTypeHurto CrimeType = "HURTO"
	// CrimeTypeVandalismo representa un acto de vandalismo
	CrimeTypeVandalismo CrimeType = "VANDALISMO"
	// CrimeTypeOtro representa otros tipos de delitos
	CrimeTypeOtro CrimeType = "OTRO"
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

// Validate valida que el delito tenga todos los campos requeridos
func (c *Crime) Validate() error {
	if c.Title == "" {
		return errors.New("el título es requerido")
	}
	if c.Description == "" {
		return errors.New("la descripción es requerida")
	}
	if c.Type == "" {
		return errors.New("el tipo de delito es requerido")
	}
	if !isValidCrimeType(c.Type) {
		return errors.New("tipo de delito inválido")
	}
	if c.Status == "" {
		return errors.New("el estado es requerido")
	}
	if !isValidCrimeStatus(c.Status) {
		return errors.New("estado inválido")
	}
	if err := c.Location.Validate(); err != nil {
		return err
	}
	if c.UpdatedAt.Before(c.CreatedAt) {
		return errors.New("la fecha de actualización no puede ser anterior a la fecha de creación")
	}
	return nil
}

// isValidCrimeType verifica si el tipo de delito es válido
func isValidCrimeType(crimeType string) bool {
	switch CrimeType(crimeType) {
	case CrimeTypeRobo, CrimeTypeAgresion, CrimeTypeHurto, CrimeTypeVandalismo, CrimeTypeOtro:
		return true
	default:
		return false
	}
}

// isValidCrimeStatus verifica si el estado del delito es válido
func isValidCrimeStatus(status string) bool {
	switch CrimeStatus(status) {
	case CrimeStatusActive, CrimeStatusInactive, CrimeStatusDeleted:
		return true
	default:
		return false
	}
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

// Validate valida que la ubicación tenga todos los campos requeridos
func (l *Location) Validate() error {
	if l.Latitude == 0 {
		return errors.New("la latitud es requerida")
	}
	if l.Longitude == 0 {
		return errors.New("la longitud es requerida")
	}
	if l.Latitude < -90 || l.Latitude > 90 {
		return errors.New("la latitud debe estar entre -90 y 90")
	}
	if l.Longitude < -180 || l.Longitude > 180 {
		return errors.New("la longitud debe estar entre -180 y 180")
	}
	if l.Address == "" {
		return errors.New("la dirección es requerida")
	}
	return nil
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
