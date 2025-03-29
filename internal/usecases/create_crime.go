package usecases

import (
	"context"
	"fmt"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
	"go-crime_map_backend/internal/domain/usecases"
)

// CreateCrimeInput representa los datos necesarios para crear un delito
type CreateCrimeInput struct {
	Title         string
	Description   string
	Type          string
	Latitude      float64
	Longitude     float64
	Address       string
	AddressNumber string
	City          string
	Province      string
	Country       string
	ZipCode       string
}

// CreateCrimeUseCase implementa la lógica de negocio para crear delitos
type CreateCrimeUseCase struct {
	crimeRepository repositories.CrimeRepository
}

// NewCreateCrimeUseCase crea una nueva instancia del caso de uso
func NewCreateCrimeUseCase(repo repositories.CrimeRepository) *CreateCrimeUseCase {
	return &CreateCrimeUseCase{
		crimeRepository: repo,
	}
}

// Execute ejecuta el caso de uso
func (uc *CreateCrimeUseCase) Execute(ctx context.Context, input usecases.CreateCrimeInput) (*entities.Crime, error) {
	// Validar datos de entrada
	if input.Title == "" {
		return nil, fmt.Errorf("el título es requerido")
	}
	if input.Description == "" {
		return nil, fmt.Errorf("la descripción es requerida")
	}
	if input.Type == "" {
		return nil, fmt.Errorf("el tipo es requerido")
	}
	if input.Latitude < -90 || input.Latitude > 90 {
		return nil, fmt.Errorf("la latitud debe estar entre -90 y 90")
	}
	if input.Longitude < -180 || input.Longitude > 180 {
		return nil, fmt.Errorf("la longitud debe estar entre -180 y 180")
	}
	if input.Address == "" {
		return nil, fmt.Errorf("la dirección es requerida")
	}
	if input.AddressNumber == "" {
		return nil, fmt.Errorf("el número de la dirección es requerido")
	}
	if input.City == "" {
		return nil, fmt.Errorf("la ciudad es requerida")
	}
	if input.Province == "" {
		return nil, fmt.Errorf("la provincia es requerida")
	}
	if input.Country == "" {
		return nil, fmt.Errorf("el país es requerido")
	}
	if input.ZipCode == "" {
		return nil, fmt.Errorf("el código postal es requerido")
	}

	// Crear entidad Crime
	crime := &entities.Crime{
		Title:       input.Title,
		Description: input.Description,
		Type:        input.Type,
		Status:      "ACTIVE",
		Location: entities.Location{
			Latitude:      input.Latitude,
			Longitude:     input.Longitude,
			Address:       input.Address,
			AddressNumber: input.AddressNumber,
			City:          input.City,
			Province:      input.Province,
			Country:       input.Country,
			ZipCode:       input.ZipCode,
		},
	}

	// Guardar en el repositorio
	if err := uc.crimeRepository.Create(ctx, crime); err != nil {
		return nil, fmt.Errorf("error al crear el delito: %w", err)
	}

	return crime, nil
}
