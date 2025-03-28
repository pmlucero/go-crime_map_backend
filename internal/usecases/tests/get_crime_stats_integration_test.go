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

func TestGetCrimeStatsUseCase_Integration(t *testing.T) {
	// Configurar base de datos de test
	db := database.SetupTestDB(t)
	defer database.CleanupTestDB(t)

	// Crear el repositorio y los casos de uso
	repo := infraRepo.NewPostgresCrimeRepository(db)
	useCase := usecases.NewGetCrimeStatsUseCase(repo)
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
		{
			Title:       "Robo de auto",
			Type:        "ROBO",
			Description: "Robo de vehículo estacionado",
			Latitude:    -34.606722,
			Longitude:   -58.384592,
			Address:     "Av. Corrientes 4567",
		},
	}

	for _, crimeInput := range crimes {
		crime, err := createCrimeUseCase.Execute(context.Background(), crimeInput)
		assert.NoError(t, err)
		assert.NotNil(t, crime)
	}

	// Ejecutar el caso de uso
	stats, err := useCase.Execute(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, stats)

	// Verificar las estadísticas
	assert.Equal(t, int64(4), stats.TotalCrimes)
	assert.Equal(t, int64(4), stats.ActiveCrimes)
	assert.Equal(t, int64(0), stats.InactiveCrimes)

	// Verificar conteo por tipo
	assert.Equal(t, int64(2), stats.CrimesByType["ROBO"])
	assert.Equal(t, int64(1), stats.CrimesByType["VANDALISMO"])
	assert.Equal(t, int64(1), stats.CrimesByType["HURTO"])

	// Verificar conteo por dirección
	assert.Equal(t, int64(1), stats.CrimesByAddress["Av. Corrientes 1234"])
	assert.Equal(t, int64(1), stats.CrimesByAddress["Av. Corrientes 2345"])
	assert.Equal(t, int64(1), stats.CrimesByAddress["Av. Corrientes 3456"])
	assert.Equal(t, int64(1), stats.CrimesByAddress["Av. Corrientes 4567"])

	// Verificar que la última actualización está presente
	assert.NotZero(t, stats.LastUpdate)
}
