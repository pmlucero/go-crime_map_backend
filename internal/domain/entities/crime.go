package entities

import (
	"errors"
	"regexp"
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
	// CrimeStatusResolved indica que el delito ha sido resuelto
	CrimeStatusResolved CrimeStatus = "RESOLVED"
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

// Crime representa un delito en el sistema
type Crime struct {
	ID          int64      `json:"id" db:"id"`
	UUID        string     `json:"uuid" db:"uuid"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Type        string     `json:"type" db:"type"`
	Status      string     `json:"status" db:"status"`
	Location    Location   `json:"location" db:",inline"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Location representa la ubicación geográfica de un delito
type Location struct {
	Latitude      float64 `json:"latitude" db:"latitude"`
	Longitude     float64 `json:"longitude" db:"longitude"`
	Address       string  `json:"address" db:"address"`
	AddressNumber string  `json:"address_number" db:"address_number"`
	City          string  `json:"city" db:"city"`
	Province      string  `json:"province" db:"province"`
	Country       string  `json:"country" db:"country"`
	ZipCode       string  `json:"zip_code" db:"zip_code"`
}

// Validate valida que el delito tenga todos los campos requeridos
func (c *Crime) Validate() error {
	if c.UUID == "" {
		return errors.New("el UUID es requerido")
	}
	if !isValidUUID(c.UUID) {
		return errors.New("UUID inválido")
	}
	if c.Type == "" {
		return errors.New("el tipo de delito es requerido")
	}
	if !isValidCrimeType(c.Type) {
		return errors.New("tipo de delito inválido")
	}
	if c.Description == "" {
		return errors.New("la descripción es requerida")
	}
	if c.Location.Latitude < -90 || c.Location.Latitude > 90 {
		return errors.New("la latitud debe estar entre -90 y 90")
	}
	if c.Location.Longitude < -180 || c.Location.Longitude > 180 {
		return errors.New("la longitud debe estar entre -180 y 180")
	}
	if c.Status == "" {
		return errors.New("el estado es requerido")
	}
	if !isValidCrimeStatus(c.Status) {
		return errors.New("estado inválido")
	}
	if c.UpdatedAt.Before(c.CreatedAt) {
		return errors.New("la fecha de actualización no puede ser anterior a la fecha de creación")
	}
	return nil
}

// isValidUUID verifica si el string es un UUID válido
func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
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
	case CrimeStatusActive, CrimeStatusInactive, CrimeStatusDeleted, CrimeStatusResolved:
		return true
	default:
		return false
	}
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

// IsValid verifica si el delito tiene todos los campos requeridos
func (c *Crime) IsValid() bool {
	return c.Type != "" &&
		c.Description != "" &&
		c.Location.Latitude >= -90 && c.Location.Latitude <= 90 &&
		c.Location.Longitude >= -180 && c.Location.Longitude <= 180
}
