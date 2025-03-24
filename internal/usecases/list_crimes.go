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
	Latitude  float64   // Latitud del centro
	Longitude float64   // Longitud del centro
	Radius    float64   // Radio en kilómetros
	Limit     int       // Límite de resultados
	Offset    int       // Desplazamiento para paginación
}

// ListCrimesUseCase maneja la lógica de negocio para listar delitos
type ListCrimesUseCase struct {
	crimeRepo repositories.CrimeRepository
}

// NewListCrimesUseCase crea una nueva instancia del caso de uso
func NewListCrimesUseCase(repo repositories.CrimeRepository) *ListCrimesUseCase {
	return &ListCrimesUseCase{
		crimeRepo: repo,
	}
}

// Execute ejecuta el caso de uso para listar delitos
func (uc *ListCrimesUseCase) Execute(ctx context.Context, input ListCrimesInput) ([]*entities.Crime, error) {
	// Validar que las fechas sean coherentes
	if !input.StartDate.IsZero() && !input.EndDate.IsZero() && input.StartDate.After(input.EndDate) {
		return nil, ErrInvalidDateRange
	}

	// Validar que el radio sea positivo si se especifican coordenadas
	if (input.Latitude != 0 || input.Longitude != 0) && input.Radius <= 0 {
		return nil, ErrInvalidRadius
	}

	// Validar que el límite sea positivo
	if input.Limit <= 0 {
		input.Limit = 10 // Valor por defecto
	}

	// Validar que el offset sea no negativo
	if input.Offset < 0 {
		input.Offset = 0
	}

	// Listar delitos con los filtros especificados
	crimes, err := uc.crimeRepo.List(ctx, repositories.ListCrimesFilter{
		Type:      input.Type,
		Status:    input.Status,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		Radius:    input.Radius,
		Limit:     input.Limit,
		Offset:    input.Offset,
	})
	if err != nil {
		return nil, err
	}

	return crimes, nil
}
