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

func strPtr(s string) *string {
	return &s
}

func TestListCrimesUseCase_Integration(t *testing.T) {
	// Configurar base de datos de test
	db := database.SetupTestDB(t)
	defer database.CleanupTestDB(t)

	// Crear el repositorio y los casos de uso
	repo := infraRepo.NewPostgresCrimeRepository(db)
	useCase := usecases.NewListCrimesUseCase(repo)
	createCrimeUseCase := usecases.NewCreateCrimeUseCase(repo)

	// Crear varios delitos de prueba
	crimes := []domain_usecases.CreateCrimeInput{
		{
			Title:       "Robo a mano armada",
			Type:        "ROBO",
			Description: "Robo a mano armada en comercio",
			Latitude:    -34.603722,
			Longitude:   -58.381592,
			Address:     "Av. Corrientes 1234",
		},
		{
			Title:       "Vandalismo",
			Type:        "VANDALISMO",
			Description: "Daños a propiedad pública",
			Latitude:    -34.604722,
			Longitude:   -58.382592,
			Address:     "Av. Corrientes 2345",
		},
		{
			Title:       "Hurto",
			Type:        "HURTO",
			Description: "Hurto de celular",
			Latitude:    -34.605722,
			Longitude:   -58.383592,
			Address:     "Av. Corrientes 3456",
		},
	}

	for _, crimeInput := range crimes {
		crime, err := createCrimeUseCase.Execute(context.Background(), crimeInput)
		assert.NoError(t, err)
		assert.NotNil(t, crime)
	}

	tests := []struct {
		name          string
		params        domain_usecases.ListCrimesParams
		expectedCount int
	}{
		{
			name: "listar todos los delitos",
			params: domain_usecases.ListCrimesParams{
				Page:  1,
				Limit: 10,
			},
			expectedCount: 3,
		},
		{
			name: "listar delitos con límite 2",
			params: domain_usecases.ListCrimesParams{
				Page:  1,
				Limit: 2,
			},
			expectedCount: 2,
		},
		{
			name: "listar delitos segunda página",
			params: domain_usecases.ListCrimesParams{
				Page:  2,
				Limit: 2,
			},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := useCase.Execute(context.Background(), tt.params)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result.Items, tt.expectedCount)

			for _, crime := range result.Items {
				assert.NotEmpty(t, crime.ID)
				assert.NotEmpty(t, crime.Title)
				assert.NotEmpty(t, crime.Type)
				assert.NotEmpty(t, crime.Description)
				assert.NotZero(t, crime.Location.Latitude)
				assert.NotZero(t, crime.Location.Longitude)
				assert.NotEmpty(t, crime.Location.Address)
				assert.NotZero(t, crime.CreatedAt)
				assert.NotZero(t, crime.UpdatedAt)
			}
		})
	}
}
