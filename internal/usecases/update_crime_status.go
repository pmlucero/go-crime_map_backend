package usecases

import (
	"context"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
)

// UpdateCrimeStatusInput representa los datos necesarios para actualizar el estado
type UpdateCrimeStatusInput struct {
	CrimeID   string
	NewStatus entities.CrimeStatus
}

// UpdateCrimeStatusUseCase maneja la lógica de negocio para actualizar el estado de un delito
type UpdateCrimeStatusUseCase struct {
	repo repositories.CrimeRepository
}

// NewUpdateCrimeStatusUseCase crea una nueva instancia del caso de uso
func NewUpdateCrimeStatusUseCase(repo repositories.CrimeRepository) *UpdateCrimeStatusUseCase {
	return &UpdateCrimeStatusUseCase{
		repo: repo,
	}
}

// Execute ejecuta el caso de uso para actualizar el estado de un delito
func (uc *UpdateCrimeStatusUseCase) Execute(ctx context.Context, input UpdateCrimeStatusInput) error {
	// Obtener el delito
	crime, err := uc.repo.GetByID(ctx, input.CrimeID)
	if err != nil {
		return err
	}

	// Verificar si el delito ya está eliminado
	if crime.Status == entities.CrimeStatusDeleted {
		return ErrCrimeAlreadyDeleted
	}

	// Validar la transición de estado
	if !isValidStatusTransition(crime.Status, input.NewStatus) {
		return ErrInvalidStatusTransition
	}

	// Actualizar el estado
	crime.Status = input.NewStatus
	return uc.repo.Update(ctx, crime)
}

// isValidStatusTransition valida si la transición de estado es válida
func isValidStatusTransition(current, new entities.CrimeStatus) bool {
	switch current {
	case entities.CrimeStatusActive:
		return new == entities.CrimeStatusInactive || new == entities.CrimeStatusDeleted
	case entities.CrimeStatusInactive:
		return new == entities.CrimeStatusActive || new == entities.CrimeStatusDeleted
	case entities.CrimeStatusDeleted:
		return false
	default:
		return false
	}
}
