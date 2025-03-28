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
	Offset    int       // Desplazamiento para paginación
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
	List(ctx context.Context, page, limit int) ([]entities.Crime, int64, error)

	// GetByID obtiene un delito por su ID
	GetByID(ctx context.Context, id string) (*entities.Crime, error)

	// Update actualiza un delito existente
	Update(ctx context.Context, crime *entities.Crime) error

	// Delete realiza una eliminación lógica de un delito
	Delete(ctx context.Context, id string) error

	// GetStats obtiene estadísticas de delitos
	GetStats(ctx context.Context) (*entities.CrimeStats, error)
}
