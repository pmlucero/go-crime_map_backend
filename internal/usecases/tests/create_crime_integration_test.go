package tests

import (
	"context"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/infrastructure/database"
	"go-crime_map_backend/internal/infrastructure/repositories"
	"go-crime_map_backend/internal/usecases"

	"github.com/stretchr/testify/assert"
)

func TestCreateCrimeUseCase_Integration(t *testing.T) {
	// Configurar base de datos de test
	db := database.SetupTestDB(t)
	defer database.CleanupTestDB(t)

	// Crear el caso de uso con el repositorio real
	repo := repositories.NewPostgresCrimeRepository(db)
	useCase := usecases.NewCreateCrimeUseCase(repo)

	// Datos de prueba
	validLocation := usecases.Location{
		Latitude:  -34.603722,
		Longitude: -58.381592,
		Address:   "Av. Corrientes 1234",
	}

	// Fecha de prueba en UTC
	testDate := time.Now().Add(-1 * time.Hour).UTC()

	tests := []struct {
		name          string
		input         usecases.CreateCrimeInput
		expectedError string
	}{
		{
			name: "creación exitosa de delito",
			input: usecases.CreateCrimeInput{
				Type:        "robo",
				Description: "Robo a mano armada",
				Location:    validLocation,
				Date:        testDate,
			},
		},
		{
			name: "error - tipo de delito vacío",
			input: usecases.CreateCrimeInput{
				Type:        "",
				Description: "Robo a mano armada",
				Location:    validLocation,
				Date:        testDate,
			},
			expectedError: "el tipo de delito es requerido",
		},
		{
			name: "error - descripción vacía",
			input: usecases.CreateCrimeInput{
				Type:        "robo",
				Description: "",
				Location:    validLocation,
				Date:        testDate,
			},
			expectedError: "la descripción es requerida",
		},
		{
			name: "error - fecha futura",
			input: usecases.CreateCrimeInput{
				Type:        "robo",
				Description: "Robo a mano armada",
				Location:    validLocation,
				Date:        time.Now().Add(24 * time.Hour).UTC(),
			},
			expectedError: "la fecha del delito no puede ser futura",
		},
		{
			name: "error - latitud inválida",
			input: usecases.CreateCrimeInput{
				Type:        "robo",
				Description: "Robo a mano armada",
				Location: usecases.Location{
					Latitude:  91.0, // Latitud inválida
					Longitude: -58.381592,
					Address:   "Av. Corrientes 1234",
				},
				Date: testDate,
			},
			expectedError: "latitud inválida",
		},
		{
			name: "error - longitud inválida",
			input: usecases.CreateCrimeInput{
				Type:        "robo",
				Description: "Robo a mano armada",
				Location: usecases.Location{
					Latitude:  -34.603722,
					Longitude: 181.0, // Longitud inválida
					Address:   "Av. Corrientes 1234",
				},
				Date: testDate,
			},
			expectedError: "longitud inválida",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ejecutar el caso de uso
			result, err := useCase.Execute(context.Background(), tt.input)

			// Verificar el resultado
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.ID)
			assert.Equal(t, tt.input.Type, result.Type)
			assert.Equal(t, tt.input.Description, result.Description)
			assert.Equal(t, tt.input.Location.Latitude, result.Location.Latitude)
			assert.Equal(t, tt.input.Location.Longitude, result.Location.Longitude)
			assert.Equal(t, tt.input.Location.Address, result.Location.Address)
			assert.Equal(t, tt.input.Date.UTC(), result.Date.UTC())
			assert.Equal(t, entities.CrimeStatusActive, result.Status)

			// Verificar que el delito se guardó en la base de datos
			savedCrime, err := repo.GetByID(context.Background(), result.ID)
			assert.NoError(t, err)
			assert.NotNil(t, savedCrime)
			assert.Equal(t, result.ID, savedCrime.ID)
			assert.Equal(t, result.Type, savedCrime.Type)
			assert.Equal(t, result.Description, savedCrime.Description)
			assert.Equal(t, result.Location.Latitude, savedCrime.Location.Latitude)
			assert.Equal(t, result.Location.Longitude, savedCrime.Location.Longitude)
			assert.Equal(t, result.Location.Address, savedCrime.Location.Address)
			assert.Equal(t, result.Date.UTC(), savedCrime.Date.UTC())
			assert.Equal(t, result.Status, savedCrime.Status)
		})
	}
}
