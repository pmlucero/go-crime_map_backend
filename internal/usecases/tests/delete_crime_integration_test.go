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

func TestDeleteCrimeUseCase_Integration(t *testing.T) {
	// Configurar base de datos de test
	db := database.SetupTestDB(t)
	defer database.CleanupTestDB(t)

	// Crear el repositorio y el caso de uso
	repo := infraRepo.NewPostgresCrimeRepository(db)
	useCase := usecases.NewDeleteCrimeUseCase(repo)

	// Crear un delito de prueba
	createCrimeUseCase := usecases.NewCreateCrimeUseCase(repo)
	input := domain_usecases.CreateCrimeInput{
		Title:       "Robo a mano armada",
		Type:        "ROBO",
		Description: "Robo a mano armada en comercio",
		Latitude:    -34.603722,
		Longitude:   -58.381592,
		Address:     "Av. Corrientes 1234",
	}
	crime, err := createCrimeUseCase.Execute(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, crime)

	tests := []struct {
		name          string
		crimeID       string
		expectedError string
	}{
		{
			name:    "eliminar delito existente",
			crimeID: crime.ID,
		},
		{
			name:          "error - delito no encontrado",
			crimeID:       "123e4567-e89b-12d3-a456-426614174000",
			expectedError: "delito no encontrado",
		},
		{
			name:          "error - delito ya eliminado",
			crimeID:       crime.ID,
			expectedError: "el delito ya ha sido eliminado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := useCase.Execute(context.Background(), tt.crimeID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}
