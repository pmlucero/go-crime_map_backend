package tests

import (
	"context"
	"testing"

	"go-crime_map_backend/internal/domain/entities"
	domain_usecases "go-crime_map_backend/internal/domain/usecases"
	"go-crime_map_backend/internal/infrastructure/database"
	infraRepo "go-crime_map_backend/internal/infrastructure/repositories"
	"go-crime_map_backend/internal/usecases"
	"go-crime_map_backend/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestUpdateCrimeStatusUseCase_Integration(t *testing.T) {
	// Configurar base de datos de test
	db := database.SetupTestDB(t)
	defer database.CleanupTestDB(t)

	// Crear el repositorio y los casos de uso
	repo := infraRepo.NewPostgresCrimeRepository(db)
	useCase := usecases.NewUpdateCrimeStatusUseCase(repo)
	createCrimeUseCase := usecases.NewCreateCrimeUseCase(repo)
	getCrimeUseCase := usecases.NewGetCrimeUseCase(repo)

	// Crear un delito de prueba
	crimeInput := domain_usecases.CreateCrimeInput{
		Title:         "Robo a mano armada",
		Type:          "ROBO",
		Description:   "Robo a mano armada en comercio",
		Latitude:      -34.603722,
		Longitude:     -58.381592,
		Address:       "Av. Corrientes",
		AddressNumber: utils.StringPtr("1234"),
		City:          utils.StringPtr("Buenos Aires"),
		Province:      utils.StringPtr("Buenos Aires"),
		Country:       utils.StringPtr("Argentina"),
		ZipCode:       utils.StringPtr("C1043"),
	}

	crime, err := createCrimeUseCase.Execute(context.Background(), crimeInput)
	assert.NoError(t, err)
	assert.NotNil(t, crime)

	tests := []struct {
		name          string
		input         domain_usecases.UpdateCrimeStatusInput
		expectedError string
	}{
		{
			name: "actualizar estado a INACTIVE",
			input: domain_usecases.UpdateCrimeStatusInput{
				ID:     crime.ID,
				Status: string(entities.CrimeStatusInactive),
			},
		},
		{
			name: "error - delito no encontrado",
			input: domain_usecases.UpdateCrimeStatusInput{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Status: string(entities.CrimeStatusInactive),
			},
			expectedError: "error al obtener el delito: error al obtener el delito: sql: no rows in result set",
		},
		{
			name: "error - estado inválido",
			input: domain_usecases.UpdateCrimeStatusInput{
				ID:     crime.ID,
				Status: "ESTADO_INVALIDO",
			},
			expectedError: "la transición de estado no es válida",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := useCase.Execute(context.Background(), tt.input)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				return
			}

			assert.NoError(t, err)

			// Verificar que el estado se actualizó correctamente
			updatedCrime, err := getCrimeUseCase.Execute(context.Background(), tt.input.ID)
			assert.NoError(t, err)
			assert.NotNil(t, updatedCrime)
			assert.Equal(t, tt.input.Status, updatedCrime.Status)
		})
	}
}
