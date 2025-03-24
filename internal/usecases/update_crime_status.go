package usecases

import (
	"context"
	"errors"
	"go-crime_map_backend/internal/domain/repositories"
	"time"
)

// UpdateCrimeStatusInput representa los datos necesarios para actualizar el estado
type UpdateCrimeStatusInput struct {
	ID        string
	NewStatus string
}

// UpdateCrimeStatusUseCase maneja la lógica de negocio para actualizar el estado de un delito
type UpdateCrimeStatusUseCase struct {
	crimeRepo repositories.CrimeRepository
}

// NewUpdateCrimeStatusUseCase crea una nueva instancia del caso de uso
func NewUpdateCrimeStatusUseCase(repo repositories.CrimeRepository) *UpdateCrimeStatusUseCase {
	return &UpdateCrimeStatusUseCase{
		crimeRepo: repo,
	}
}

// Execute ejecuta el caso de uso para actualizar el estado de un delito
func (uc *UpdateCrimeStatusUseCase) Execute(ctx context.Context, input UpdateCrimeStatusInput) error {
	crime, err := uc.crimeRepo.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}

	if crime == nil {
		return errors.New("delito no encontrado")
	}

	if crime.DeletedAt != nil {
		return errors.New("no se puede actualizar un delito eliminado")
	}

	if !isValidStatusTransition(crime.Status, input.NewStatus) {
		return errors.New("transición de estado no válida")
	}

	crime.Status = input.NewStatus
	crime.UpdatedAt = time.Now()

	return uc.crimeRepo.Update(ctx, crime)
}

// isValidStatusTransition valida si la transición de estado es válida
func isValidStatusTransition(currentStatus, newStatus string) bool {
	validTransitions := map[string][]string{
		"active":   {"inactive"},
		"inactive": {"active"},
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			return true
		}
	}

	return false
}
