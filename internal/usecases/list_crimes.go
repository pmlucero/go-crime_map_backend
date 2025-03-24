package usecases

import (
	"context"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
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

// ListCrimesUseCase maneja la lógica de negocio para listar delitos
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

// Execute ejecuta el caso de uso para listar delitos
func (uc *ListCrimesUseCase) Execute(ctx context.Context, params ListCrimesParams) (*entities.CrimeList, error) {
	// Ajustar las fechas para incluir todo el día
	var startDate, endDate *time.Time
	if params.StartDate != nil {
		start := time.Date(params.StartDate.Year(), params.StartDate.Month(), params.StartDate.Day(), 0, 0, 0, 0, params.StartDate.Location())
		startDate = &start
	}
	if params.EndDate != nil {
		end := time.Date(params.EndDate.Year(), params.EndDate.Month(), params.EndDate.Day(), 23, 59, 59, 999999999, params.EndDate.Location())
		endDate = &end
	}

	// Validar y ajustar parámetros de paginación
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
	crimes, total, err := uc.crimeRepository.List(ctx, params.Page, params.Limit, startDate, endDate, params.Type, params.Status)
	if err != nil {
		return nil, err
	}

	// Calcular total de páginas
	totalPages := (total + int64(params.Limit) - 1) / int64(params.Limit)

	return &entities.CrimeList{
		Crimes:      crimes,
		Total:       total,
		Page:        params.Page,
		Limit:       params.Limit,
		TotalPages:  int(totalPages),
		HasNextPage: params.Page < int(totalPages),
	}, nil
}
