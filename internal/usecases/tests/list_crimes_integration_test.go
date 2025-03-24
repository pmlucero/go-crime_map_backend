package tests

import (
	"context"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/infrastructure/database"
	infraRepo "go-crime_map_backend/internal/infrastructure/repositories"
	"go-crime_map_backend/internal/usecases"

	"github.com/stretchr/testify/assert"
)

func TestListCrimesUseCase_Integration(t *testing.T) {
	// Configurar base de datos de test
	db := database.SetupTestDB(t)
	defer database.CleanupTestDB(t)

	// Crear el repositorio y los casos de uso
	repo := infraRepo.NewPostgresCrimeRepository(db)
	listUseCase := usecases.NewListCrimesUseCase(repo)
	createUseCase := usecases.NewCreateCrimeUseCase(repo)
	updateStatusUseCase := usecases.NewUpdateCrimeStatusUseCase(repo)

	// Crear delitos de prueba
	crimes := []struct {
		Type        string
		Description string
		Status      entities.CrimeStatus
		Date        time.Time
		Latitude    float64
		Longitude   float64
	}{
		{
			Type:        "robo",
			Description: "Descripción del robo 1",
			Status:      entities.CrimeStatusActive,
			Date:        time.Now().Add(-24 * time.Hour),
			Latitude:    -34.603722,
			Longitude:   -58.381592,
		},
		{
			Type:        "hurto",
			Description: "Descripción del hurto 1",
			Status:      entities.CrimeStatusActive,
			Date:        time.Now().Add(-48 * time.Hour),
			Latitude:    -34.588722,
			Longitude:   -58.392592,
		},
		{
			Type:        "robo",
			Description: "Descripción del robo 2",
			Status:      entities.CrimeStatusInactive,
			Date:        time.Now().Add(-72 * time.Hour),
			Latitude:    -34.613722,
			Longitude:   -58.371592,
		},
	}

	var createdCrimes []entities.Crime
	for _, c := range crimes {
		crime, err := createUseCase.Execute(context.Background(), usecases.CreateCrimeInput{
			Type:        c.Type,
			Description: c.Description,
			Date:        c.Date,
			Location: usecases.Location{
				Latitude:  c.Latitude,
				Longitude: c.Longitude,
				Address:   "Test Address",
			},
		})
		assert.NoError(t, err)
		createdCrimes = append(createdCrimes, *crime)

		if c.Status == entities.CrimeStatusInactive {
			err = updateStatusUseCase.Execute(context.Background(), usecases.UpdateCrimeStatusInput{
				CrimeID:   crime.ID,
				NewStatus: entities.CrimeStatusInactive,
			})
			assert.NoError(t, err)
		}
	}

	tests := []struct {
		name          string
		input         usecases.ListCrimesInput
		expectedTypes []string
		expectedCount int
		expectedError string
	}{
		{
			name:          "listar todos los delitos",
			input:         usecases.ListCrimesInput{},
			expectedTypes: []string{"robo", "hurto", "robo"},
			expectedCount: 3,
			expectedError: "",
		},
		{
			name: "filtrar por tipo de delito",
			input: usecases.ListCrimesInput{
				Type: "robo",
			},
			expectedTypes: []string{"robo", "robo"},
			expectedCount: 2,
			expectedError: "",
		},
		{
			name: "filtrar por estado",
			input: usecases.ListCrimesInput{
				Status: string(entities.CrimeStatusInactive),
			},
			expectedTypes: []string{"robo"},
			expectedCount: 1,
			expectedError: "",
		},
		{
			name: "filtrar por fecha",
			input: usecases.ListCrimesInput{
				StartDate: time.Now().Add(-49 * time.Hour),
				EndDate:   time.Now(),
			},
			expectedTypes: []string{"hurto", "robo"},
			expectedCount: 2,
			expectedError: "",
		},
		{
			name: "filtrar por ubicación",
			input: usecases.ListCrimesInput{
				Latitude:  -34.603722,
				Longitude: -58.381592,
				Radius:    0.5,
			},
			expectedTypes: []string{"robo"},
			expectedCount: 1,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ejecutar el caso de uso
			result, err := listUseCase.Execute(context.Background(), tt.input)

			// Verificar el resultado
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, tt.expectedCount)

			// Verificar los tipos de los delitos
			for i, crime := range result {
				assert.Equal(t, tt.expectedTypes[i], crime.Type)
			}
		})
	}
}
