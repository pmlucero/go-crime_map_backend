package usecases

import (
	"context"
	"errors"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"

	"github.com/google/uuid"
)

var (
	// ErrInvalidType se retorna cuando el tipo de delito es inválido
	ErrInvalidType = errors.New("el tipo de delito es inválido")

	// ErrEmptyDescription se retorna cuando la descripción está vacía
	ErrEmptyDescription = errors.New("la descripción es requerida")

	// ErrDescriptionTooLong se retorna cuando la descripción excede el límite de caracteres
	ErrDescriptionTooLong = errors.New("la descripción no puede exceder los 500 caracteres")

	// ErrFutureDate se retorna cuando la fecha es futura
	ErrFutureDate = errors.New("la fecha del delito no puede ser futura")

	// ErrInvalidLatitude se retorna cuando la latitud es inválida
	ErrInvalidLatitude = errors.New("latitud inválida")

	// ErrInvalidLongitude se retorna cuando la longitud es inválida
	ErrInvalidLongitude = errors.New("longitud inválida")

	// ErrDuplicateCrime se retorna cuando se intenta crear un delito duplicado
	ErrDuplicateCrime = errors.New("ya existe un delito con los mismos datos")

	// validCrimeTypes define los tipos de delito válidos
	validCrimeTypes = map[string]bool{
		"ROBO":         true,
		"HURTO":        true,
		"VANDALISMO":   true,
		"AGRESION":     true,
		"FRAUDE":       true,
		"TRAFICO":      true,
		"ACOSO":        true,
		"VIOLENCIA":    true,
		"ALLANAMIENTO": true,
		"ESTAFA":       true,
	}

	// maxDescriptionLength define la longitud máxima permitida para la descripción
	maxDescriptionLength = 500
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
		return nil, ErrFutureDate
	}

	// Validar que el tipo de delito sea válido
	if !validCrimeTypes[input.Type] {
		return nil, ErrInvalidType
	}

	// Validar que la descripción no esté vacía
	if input.Description == "" {
		return nil, ErrEmptyDescription
	}

	// Validar que la descripción no exceda el límite de caracteres
	if len(input.Description) > maxDescriptionLength {
		return nil, ErrDescriptionTooLong
	}

	// Validar que la ubicación sea válida
	if input.Location.Latitude <= -90 || input.Location.Latitude >= 90 {
		return nil, ErrInvalidLatitude
	}
	if input.Location.Longitude <= -180 || input.Location.Longitude >= 180 {
		return nil, ErrInvalidLongitude
	}

	// Verificar si existe un delito duplicado
	crimes, err := uc.crimeRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, crime := range crimes {
		// Comparar fechas con una tolerancia de 1 minuto
		dateDiff := crime.Date.Sub(input.Date)
		if dateDiff < 0 {
			dateDiff = -dateDiff
		}
		if dateDiff > time.Minute {
			continue
		}

		// Comparar coordenadas con una tolerancia de 0.000001 grados
		latDiff := crime.Location.Latitude - input.Location.Latitude
		if latDiff < 0 {
			latDiff = -latDiff
		}
		if latDiff > 0.000001 {
			continue
		}

		lonDiff := crime.Location.Longitude - input.Location.Longitude
		if lonDiff < 0 {
			lonDiff = -lonDiff
		}
		if lonDiff > 0.000001 {
			continue
		}

		// Si llegamos aquí, todas las comparaciones pasaron
		if crime.Type == input.Type &&
			crime.Description == input.Description {
			return nil, ErrDuplicateCrime
		}
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Guardar en el repositorio
	if err := uc.crimeRepo.Create(ctx, crime); err != nil {
		return nil, err
	}

	return crime, nil
}

// generateID genera un ID único para el delito usando UUID v4
func generateID() string {
	return uuid.New().String()
}
