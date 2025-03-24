package tests

import (
	"context"
	"strings"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
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
	testDate := time.Now().Add(-1 * time.Hour).UTC()
	createCrimeUseCase := usecases.NewCreateCrimeUseCase(repo)
	input := usecases.CreateCrimeInput{
		Type:        "robo",
		Description: "Robo a mano armada",
		Location: usecases.Location{
			Latitude:  -34.603722,
			Longitude: -58.381592,
			Address:   "Av. Corrientes 1234",
		},
		Date: testDate,
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
			expectedError: "recurso no encontrado",
		},
		{
			name:          "error - delito ya eliminado",
			crimeID:       crime.ID,
			expectedError: "el delito ya fue eliminado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ejecutar el caso de uso
			err := useCase.Execute(context.Background(), tt.crimeID)

			// Verificar el resultado
			if tt.expectedError != "" {
				assert.Error(t, err)
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("%s does not contain %s", err.Error(), tt.expectedError)
				}
				return
			}

			assert.NoError(t, err)

			// Verificar que el delito se marcó como eliminado en la base de datos
			deletedCrime, err := repo.GetByID(context.Background(), tt.crimeID)
			assert.NoError(t, err)
			assert.NotNil(t, deletedCrime)
			assert.Equal(t, entities.CrimeStatusDeleted, deletedCrime.Status)
		})
	}
}
