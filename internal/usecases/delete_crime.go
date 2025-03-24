package usecases

import (
	"context"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
)

type DeleteCrimeUseCase struct {
	repo repositories.CrimeRepository
}

func NewDeleteCrimeUseCase(repo repositories.CrimeRepository) *DeleteCrimeUseCase {
	return &DeleteCrimeUseCase{
		repo: repo,
	}
}

func (uc *DeleteCrimeUseCase) Execute(ctx context.Context, crimeID string) error {
	// Obtener el delito
	crime, err := uc.repo.GetByID(ctx, crimeID)
	if err != nil {
		return err
	}

	// Verificar si el delito ya está eliminado
	if crime.Status == entities.CrimeStatusDeleted {
		return ErrCrimeAlreadyDeleted
	}

	// Actualizar el estado a eliminado
	crime.Status = entities.CrimeStatusDeleted
	return uc.repo.Update(ctx, crime)
}
