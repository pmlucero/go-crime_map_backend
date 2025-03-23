package entities

import "time"

// Crime representa un delito reportado en el sistema
type Crime struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`        // Tipo de delito (robo, asalto, etc.)
	Description string    `json:"description"` // Descripción detallada del delito
	Location    Location  `json:"location"`    // Ubicación donde ocurrió el delito
	Date        time.Time `json:"date"`        // Fecha y hora del delito
	CreatedAt   time.Time `json:"created_at"`  // Fecha de creación del registro
	UpdatedAt   time.Time `json:"updated_at"`  // Fecha de última actualización
}

// Location representa la ubicación geográfica de un delito
type Location struct {
	Latitude  float64 `json:"latitude"`  // Latitud
	Longitude float64 `json:"longitude"` // Longitud
	Address   string  `json:"address"`   // Dirección descriptiva
}
