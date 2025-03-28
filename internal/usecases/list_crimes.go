package usecases

import (
	"context"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
	"go-crime_map_backend/internal/domain/usecases"
)

// ListCrimesInput representa los filtros para listar delitos
type ListCrimesInput struct {
	Type      string    // Tipo de delito
	Status    string    // Estado del delito
	StartDate time.Time // Fecha inicial
	EndDate   time.Time // Fecha final
	Limit     int       // Límite de resultados por página
	Offset    int       // Desplazamiento para paginación
}

// ListCrimesOutput representa el resultado de listar delitos
type ListCrimesOutput struct {
	Crimes     []entities.Crime // Lista de delitos
	TotalCount int64            // Total de delitos que coinciden con los filtros
	Page       int              // Página actual
	TotalPages int              // Total de páginas
}

// ListCrimesUseCase implementa la lógica de negocio para listar delitos
type ListCrimesUseCase struct {
	crimeRepository repositories.CrimeRepository
}

// ListCrimesParams representa los parámetros para listar delitos
type ListCrimesParams struct {
	Page      int
	Limit     int
	StartDate *time.Time
	EndDate   *time.Time
	Type      *string
	Status    *string
}

// NewListCrimesUseCase crea una nueva instancia del caso de uso
func NewListCrimesUseCase(repo repositories.CrimeRepository) *ListCrimesUseCase {
	return &ListCrimesUseCase{
		crimeRepository: repo,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListCrimesUseCase) Execute(ctx context.Context, params usecases.ListCrimesParams) (*entities.CrimeList, error) {
	// Validar parámetros
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	// Obtener delitos del repositorio
	crimes, total, err := uc.crimeRepository.List(ctx, params.Page, params.Limit)
	if err != nil {
		return nil, err
	}

	return &entities.CrimeList{
		Items: crimes,
		Total: total,
	}, nil
}
