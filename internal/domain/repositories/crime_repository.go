package repositories

import (
	"context"
	"time"

	"go-crime_map_backend/internal/domain/entities"
)

// ListCrimesFilter define los filtros para listar delitos
type ListCrimesFilter struct {
	Type      string    // Tipo de delito
	Status    string    // Estado del delito
	StartDate time.Time // Fecha inicial
	EndDate   time.Time // Fecha final
	Limit     int       // Límite de resultados por página
	Page      int       // Número de página
	Latitude  float64   // Latitud para búsqueda por ubicación
	Longitude float64   // Longitud para búsqueda por ubicación
	RadiusKm  float64   // Radio de búsqueda en kilómetros
}

// ListCrimesResult representa el resultado de listar delitos
type ListCrimesResult struct {
	Crimes     []*entities.Crime // Lista de delitos
	TotalCount int64             // Total de delitos que coinciden con los filtros
}

// CrimeRepository define la interfaz para el repositorio de delitos
type CrimeRepository interface {
	// Create crea un nuevo delito
	Create(ctx context.Context, crime *entities.Crime) error

	// List obtiene una lista paginada de delitos
	List(ctx context.Context, page, limit int, startDate, endDate *time.Time, crimeType, status *string) ([]entities.Crime, int64, error)

	// GetByID obtiene un delito por su ID numérico
	GetByID(ctx context.Context, id int64) (*entities.Crime, error)

	// GetByUUID obtiene un delito por su UUID
	GetByUUID(ctx context.Context, uuid string) (*entities.Crime, error)

	// Update actualiza un delito existente
	Update(ctx context.Context, crime *entities.Crime) error

	// Delete elimina un delito por su ID
	Delete(ctx context.Context, id int64) error

	// GetStats obtiene estadísticas de delitos
	GetStats(ctx context.Context) (*entities.CrimeStats, error)
}
