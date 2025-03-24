package usecases

import (
	"context"
	"errors"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"

	"github.com/google/uuid"
)

// CreateCrimeInput representa los datos necesarios para crear un delito
type CreateCrimeInput struct {
	Title       string
	Description string
	Type        string
	Location    entities.Location
}

// CreateCrimeUseCase maneja la lógica de negocio para crear un nuevo delito
type CreateCrimeUseCase struct {
	crimeRepo repositories.CrimeRepository
}

// NewCreateCrimeUseCase crea una nueva instancia del caso de uso
func NewCreateCrimeUseCase(repo repositories.CrimeRepository) *CreateCrimeUseCase {
	return &CreateCrimeUseCase{
		crimeRepo: repo,
	}
}

// Execute ejecuta el caso de uso para crear un nuevo delito
func (uc *CreateCrimeUseCase) Execute(ctx context.Context, input CreateCrimeInput) (*entities.Crime, error) {
	now := time.Now().UTC()

	// Validar que el tipo de delito no esté vacío
	if input.Type == "" {
		return nil, errors.New("el tipo de delito es requerido")
	}

	// Validar que la descripción no esté vacía
	if input.Description == "" {
		return nil, errors.New("la descripción es requerida")
	}

	// Crear la entidad Crime
	crime := &entities.Crime{
		ID:          uuid.New().String(),
		Title:       input.Title,
		Description: input.Description,
		Type:        input.Type,
		Status:      "ACTIVE",
		Location:    input.Location,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Guardar en el repositorio
	if err := uc.crimeRepo.Create(ctx, crime); err != nil {
		return nil, err
	}

	return crime, nil
}
