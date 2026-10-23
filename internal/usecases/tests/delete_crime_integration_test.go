package tests

import (
	"context"
	"testing"

	"go-crime_map_backend/internal/domain/entities"
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
	input1 := domain_usecases.CreateCrimeInput{
		Title:         "Robo de auto",
		Description:   "Me robaron el auto estacionado",
		Type:          "ROBO",
		Latitude:      -34.603722,
		Longitude:     -58.381592,
		Address:       "Av. 9 de Julio",
		AddressNumber: "1000",
		City:          "Buenos Aires",
		Province:      "Buenos Aires",
		Country:       "Argentina",
		ZipCode:       "C1043",
	}
	crime1, err := createCrimeUseCase.Execute(context.Background(), input1)
	assert.NoError(t, err)
	assert.NotNil(t, crime1)

	// Crear un segundo delito para el caso de "ya eliminado"
	input2 := domain_usecases.CreateCrimeInput{
		Title:         "Robo de moto",
		Description:   "Me robaron la moto estacionada",
		Type:          "ROBO",
		Latitude:      -34.603722,
		Longitude:     -58.381592,
		Address:       "Av. Corrientes",
		AddressNumber: "2000",
		City:          "Buenos Aires",
		Province:      "Buenos Aires",
		Country:       "Argentina",
		ZipCode:       "C1043",
	}
	crime2, err := createCrimeUseCase.Execute(context.Background(), input2)
	assert.NoError(t, err)
	assert.NotNil(t, crime2)

	tests := []struct {
		name          string
		crimeID       string
		expectedError string
	}{
		{
			name:    "eliminar delito existente",
			crimeID: crime1.UUID,
		},
		{
			name:          "error - delito no encontrado",
			crimeID:       "123e4567-e89b-12d3-a456-426614174000",
			expectedError: "error al obtener el delito: delito no encontrado con UUID 123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name:          "Error al intentar eliminar un delito que ya fue eliminado",
			crimeID:       crime2.UUID,
			expectedError: "el delito ya fue eliminado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			if tt.name == "Error al intentar eliminar un delito que ya fue eliminado" {
				err := useCase.Execute(context.Background(), tt.crimeID)
				assert.NoError(t, err)
			}
			// Act
			err := useCase.Execute(context.Background(), tt.crimeID)
			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				return
			}
			assert.NoError(t, err)
			var status string
			err = db.QueryRowContext(context.Background(), "SELECT status FROM crimes WHERE uuid = $1", tt.crimeID).Scan(&status)
			assert.NoError(t, err)
			assert.Equal(t, string(entities.CrimeStatusDeleted), status)
		})
	}
}
