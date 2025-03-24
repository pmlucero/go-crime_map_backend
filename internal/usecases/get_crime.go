package usecases

import (
	"context"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
)

// GetCrimeUseCase maneja la lógica de negocio para obtener un delito por ID
type GetCrimeUseCase struct {
	crimeRepo repositories.CrimeRepository
}

// NewGetCrimeUseCase crea una nueva instancia del caso de uso
func NewGetCrimeUseCase(repo repositories.CrimeRepository) *GetCrimeUseCase {
	return &GetCrimeUseCase{
		crimeRepo: repo,
	}
}

// Execute ejecuta el caso de uso para obtener un delito por ID
func (uc *GetCrimeUseCase) Execute(ctx context.Context, id string) (*entities.Crime, error) {
	crime, err := uc.crimeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return crime, nil
}
