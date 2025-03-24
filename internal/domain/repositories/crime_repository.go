package repositories

import (
	"context"
	"time"

	"go-crime_map_backend/internal/domain/entities"
)

// ListCrimesFilter representa los filtros para listar delitos
type ListCrimesFilter struct {
	Type      string
	Status    string
	StartDate time.Time
	EndDate   time.Time
	Latitude  float64
	Longitude float64
	Radius    float64
	Limit     int
	Offset    int
}

// CrimeRepository define la interfaz para el repositorio de delitos
type CrimeRepository interface {
	// Create crea un nuevo delito
	Create(ctx context.Context, crime *entities.Crime) error

	// GetByID obtiene un delito por su ID
	GetByID(ctx context.Context, id string) (*entities.Crime, error)

	// GetAll obtiene todos los delitos
	GetAll(ctx context.Context) ([]*entities.Crime, error)

	// Update actualiza un delito existente
	Update(ctx context.Context, crime *entities.Crime) error

	// Delete elimina un delito por su ID
	Delete(ctx context.Context, id string) error

	// List obtiene una lista de delitos con los filtros especificados
	List(ctx context.Context, filter ListCrimesFilter) ([]*entities.Crime, error)
}
