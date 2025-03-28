package usecases

import (
	"context"
	"fmt"

	"go-crime_map_backend/internal/domain/repositories"
	"go-crime_map_backend/internal/domain/usecases"
)

// UpdateCrimeStatusInput representa los datos necesarios para actualizar el estado
type UpdateCrimeStatusInput struct {
	ID     string
	Status string
}

// UpdateCrimeStatusUseCase implementa la lógica de negocio para actualizar el estado de un delito
type UpdateCrimeStatusUseCase struct {
	crimeRepository repositories.CrimeRepository
}

// NewUpdateCrimeStatusUseCase crea una nueva instancia del caso de uso
func NewUpdateCrimeStatusUseCase(repo repositories.CrimeRepository) *UpdateCrimeStatusUseCase {
	return &UpdateCrimeStatusUseCase{
		crimeRepository: repo,
	}
}

// Execute ejecuta el caso de uso
func (uc *UpdateCrimeStatusUseCase) Execute(ctx context.Context, input usecases.UpdateCrimeStatusInput) error {
	// Validar datos de entrada
	if input.ID == "" {
		return fmt.Errorf("el ID es requerido")
	}
	if input.Status == "" {
		return fmt.Errorf("el estado es requerido")
	}

	// Obtener el delito del repositorio
	crime, err := uc.crimeRepository.GetByID(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("error al obtener el delito: %w", err)
	}

	if crime == nil {
		return fmt.Errorf("delito no encontrado")
	}

	// Actualizar el estado
	crime.Status = input.Status

	// Guardar en el repositorio
	if err := uc.crimeRepository.Update(ctx, crime); err != nil {
		return fmt.Errorf("error al actualizar el delito: %w", err)
	}

	return nil
}
