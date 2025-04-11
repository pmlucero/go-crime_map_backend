package usecases

import (
	"context"
	"fmt"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
)

// DeleteCrimeUseCase implementa la lógica de negocio para eliminar delitos
type DeleteCrimeUseCase struct {
	crimeRepository repositories.CrimeRepository
}

// NewDeleteCrimeUseCase crea una nueva instancia del caso de uso
func NewDeleteCrimeUseCase(repo repositories.CrimeRepository) *DeleteCrimeUseCase {
	return &DeleteCrimeUseCase{
		crimeRepository: repo,
	}
}

// Execute ejecuta el caso de uso
func (uc *DeleteCrimeUseCase) Execute(ctx context.Context, id string) error {
	// Validar datos de entrada
	if id == "" {
		return fmt.Errorf("el ID es requerido")
	}

	// Verificar que el delito existe
	crime, err := uc.crimeRepository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error al obtener el delito: %w", err)
	}

	if crime == nil {
		return fmt.Errorf("delito no encontrado")
	}

	// Verificar que el delito no esté ya eliminado
	if crime.Status == string(entities.CrimeStatusDeleted) {
		return fmt.Errorf("el delito ya fue eliminado")
	}

	// Eliminar del repositorio
	if err := uc.crimeRepository.Delete(ctx, id); err != nil {
		return fmt.Errorf("error al eliminar el delito: %w", err)
	}

	return nil
}
