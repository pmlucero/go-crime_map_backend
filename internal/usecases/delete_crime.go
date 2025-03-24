package usecases

import (
	"context"
	"errors"
	"go-crime_map_backend/internal/domain/repositories"
	"time"
)

type DeleteCrimeUseCase struct {
	crimeRepo repositories.CrimeRepository
}

func NewDeleteCrimeUseCase(repo repositories.CrimeRepository) *DeleteCrimeUseCase {
	return &DeleteCrimeUseCase{
		crimeRepo: repo,
	}
}

func (uc *DeleteCrimeUseCase) Execute(ctx context.Context, id string) error {
	crime, err := uc.crimeRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if crime == nil {
		return errors.New("delito no encontrado")
	}

	if crime.DeletedAt != nil {
		return errors.New("el delito ya ha sido eliminado")
	}

	now := time.Now()
	crime.Status = "deleted"
	crime.DeletedAt = &now
	crime.UpdatedAt = now

	return uc.crimeRepo.Update(ctx, crime)
}
