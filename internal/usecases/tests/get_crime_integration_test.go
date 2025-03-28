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

func TestGetCrimeUseCase_Integration(t *testing.T) {
	// Configurar base de datos de test
	db := database.SetupTestDB(t)
	defer database.CleanupTestDB(t)

	// Crear el repositorio y los casos de uso
	repo := infraRepo.NewPostgresCrimeRepository(db)
	useCase := usecases.NewGetCrimeUseCase(repo)
	createCrimeUseCase := usecases.NewCreateCrimeUseCase(repo)

	// Crear un delito de prueba
	input := domain_usecases.CreateCrimeInput{
		Title:       "Robo a mano armada",
		Type:        "ROBO",
		Description: "Robo a mano armada en comercio",
		Latitude:    -34.603722,
		Longitude:   -58.381592,
		Address:     "Av. Corrientes 1234",
	}
	createdCrime, err := createCrimeUseCase.Execute(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, createdCrime)

	tests := []struct {
		name          string
		crimeID       string
		expectedError string
	}{
		{
			name:    "obtener delito existente",
			crimeID: createdCrime.ID,
		},
		{
			name:          "error - delito no encontrado",
			crimeID:       "123e4567-e89b-12d3-a456-426614174000",
			expectedError: "delito no encontrado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crime, err := useCase.Execute(context.Background(), tt.crimeID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, crime)
			assert.Equal(t, createdCrime.ID, crime.ID)
			assert.Equal(t, input.Title, crime.Title)
			assert.Equal(t, input.Type, crime.Type)
			assert.Equal(t, input.Description, crime.Description)
			assert.Equal(t, input.Latitude, crime.Location.Latitude)
			assert.Equal(t, input.Longitude, crime.Location.Longitude)
			assert.Equal(t, input.Address, crime.Location.Address)
		})
	}
}
