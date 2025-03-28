package usecases

import (
	"context"
	"fmt"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
)

// GetCrimeUseCase implementa la lógica de negocio para obtener un delito por ID
type GetCrimeUseCase struct {
	crimeRepository repositories.CrimeRepository
}

// NewGetCrimeUseCase crea una nueva instancia del caso de uso
func NewGetCrimeUseCase(repo repositories.CrimeRepository) *GetCrimeUseCase {
	return &GetCrimeUseCase{
		crimeRepository: repo,
	}
}

// Execute ejecuta el caso de uso
func (uc *GetCrimeUseCase) Execute(ctx context.Context, id string) (*entities.Crime, error) {
	// Validar que el ID no esté vacío
	if id == "" {
		return nil, fmt.Errorf("el ID del delito es requerido")
	}

	// Obtener el delito del repositorio
	crime, err := uc.crimeRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el delito: %w", err)
	}

	if crime == nil {
		return nil, fmt.Errorf("delito no encontrado")
	}

	return crime, nil
}
