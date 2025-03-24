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
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Location    Location  `json:"location"`
	Date        time.Time `json:"date"`
}

// Location representa la ubicación del delito
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
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
	// Validar que la fecha no sea futura
	if input.Date.After(time.Now()) {
		return nil, errors.New("la fecha del delito no puede ser futura")
	}

	// Validar que el tipo de delito no esté vacío
	if input.Type == "" {
		return nil, errors.New("el tipo de delito es requerido")
	}

	// Validar que la descripción no esté vacía
	if input.Description == "" {
		return nil, errors.New("la descripción es requerida")
	}

	// Validar que la ubicación sea válida
	if input.Location.Latitude < -90 || input.Location.Latitude > 90 {
		return nil, errors.New("latitud inválida")
	}
	if input.Location.Longitude < -180 || input.Location.Longitude > 180 {
		return nil, errors.New("longitud inválida")
	}

	// Crear la entidad Crime
	crime := &entities.Crime{
		ID:          generateID(),
		Type:        input.Type,
		Description: input.Description,
		Location: entities.Location{
			Latitude:  input.Location.Latitude,
			Longitude: input.Location.Longitude,
			Address:   input.Location.Address,
		},
		Date:      input.Date,
		Status:    entities.CrimeStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Guardar en el repositorio
	if err := uc.crimeRepo.Create(ctx, crime); err != nil {
		return nil, err
	}

	return crime, nil
}

// generateID genera un ID único para el delito
func generateID() string {
	return uuid.New().String()
}
