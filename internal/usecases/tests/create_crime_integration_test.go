package tests

import (
	"context"
	"testing"

	domain_usecases "go-crime_map_backend/internal/domain/usecases"
	"go-crime_map_backend/internal/infrastructure/database"
	infraRepo "go-crime_map_backend/internal/infrastructure/repositories"
	"go-crime_map_backend/internal/usecases"

	"github.com/stretchr/testify/assert"
)

func TestCreateCrimeUseCase_Integration(t *testing.T) {
	// Configurar base de datos de test
	db := database.SetupTestDB(t)
	defer database.CleanupTestDB(t)

	// Crear el repositorio y el caso de uso
	repo := infraRepo.NewPostgresCrimeRepository(db)
	useCase := usecases.NewCreateCrimeUseCase(repo)

	tests := []struct {
		name          string
		input         domain_usecases.CreateCrimeInput
		expectedError string
	}{
		{
			name: "crear delito exitosamente",
			input: domain_usecases.CreateCrimeInput{
				Title:       "Robo a mano armada",
				Description: "Robo a mano armada en comercio",
				Type:        "ROBO",
				Latitude:    -34.603722,
				Longitude:   -58.381592,
				Address:     "Av. Corrientes 1234",
			},
		},
		{
			name: "error - título vacío",
			input: domain_usecases.CreateCrimeInput{
				Description: "Robo a mano armada en comercio",
				Type:        "ROBO",
				Latitude:    -34.603722,
				Longitude:   -58.381592,
				Address:     "Av. Corrientes 1234",
			},
			expectedError: "el título es requerido",
		},
		{
			name: "error - descripción vacía",
			input: domain_usecases.CreateCrimeInput{
				Title:     "Robo a mano armada",
				Type:      "ROBO",
				Latitude:  -34.603722,
				Longitude: -58.381592,
				Address:   "Av. Corrientes 1234",
			},
			expectedError: "la descripción es requerida",
		},
		{
			name: "error - tipo vacío",
			input: domain_usecases.CreateCrimeInput{
				Title:       "Robo a mano armada",
				Description: "Robo a mano armada en comercio",
				Latitude:    -34.603722,
				Longitude:   -58.381592,
				Address:     "Av. Corrientes 1234",
			},
			expectedError: "el tipo es requerido",
		},
		{
			name: "error - latitud inválida",
			input: domain_usecases.CreateCrimeInput{
				Title:       "Robo a mano armada",
				Description: "Robo a mano armada en comercio",
				Type:        "ROBO",
				Latitude:    91,
				Longitude:   -58.381592,
				Address:     "Av. Corrientes 1234",
			},
			expectedError: "la latitud debe estar entre -90 y 90",
		},
		{
			name: "error - longitud inválida",
			input: domain_usecases.CreateCrimeInput{
				Title:       "Robo a mano armada",
				Description: "Robo a mano armada en comercio",
				Type:        "ROBO",
				Latitude:    -34.603722,
				Longitude:   181,
				Address:     "Av. Corrientes 1234",
			},
			expectedError: "la longitud debe estar entre -180 y 180",
		},
		{
			name: "error - dirección vacía",
			input: domain_usecases.CreateCrimeInput{
				Title:       "Robo a mano armada",
				Description: "Robo a mano armada en comercio",
				Type:        "ROBO",
				Latitude:    -34.603722,
				Longitude:   -58.381592,
			},
			expectedError: "la dirección es requerida",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crime, err := useCase.Execute(context.Background(), tt.input)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, crime)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, crime)
			assert.Equal(t, tt.input.Title, crime.Title)
			assert.Equal(t, tt.input.Description, crime.Description)
			assert.Equal(t, tt.input.Type, crime.Type)
			assert.Equal(t, tt.input.Latitude, crime.Location.Latitude)
			assert.Equal(t, tt.input.Longitude, crime.Location.Longitude)
			assert.Equal(t, tt.input.Address, crime.Location.Address)
		})
	}
}
