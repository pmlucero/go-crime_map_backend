package usecases

import (
	"context"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestGetCrimeStatsUseCase_Execute verifica que las estadísticas
// se obtengan correctamente del repositorio
func TestGetCrimeStatsUseCase_Execute(t *testing.T) {
	// Crear mock del repositorio
	mockRepo := new(mocks.MockCrimeRepository)
	useCase := NewGetCrimeStatsUseCase(mockRepo)

	// Crear contexto
	ctx := context.Background()

	// Crear datos de prueba
	stats := &entities.CrimeStats{
		TotalCrimes:    100,
		ActiveCrimes:   80,
		InactiveCrimes: 20,
		CrimesByType: map[string]int64{
			"ROBO":   50,
			"ASALTO": 30,
			"HURTO":  20,
		},
		CrimesByStatus: map[string]int64{
			"ACTIVE":   80,
			"INACTIVE": 20,
		},
		CrimesByLocation: map[string]int64{
			"CABA": 60,
			"GBA":  40,
		},
		CrimesByAddress: map[string]int64{
			"Av. Corrientes": 30,
			"Av. Rivadavia":  20,
			"Av. 9 de Julio": 10,
		},
		LastUpdate: time.Now(),
	}

	// Configurar el mock
	mockRepo.On("GetStats", mock.Anything).Return(stats, nil)

	// Ejecutar el caso de uso
	result, err := useCase.Execute(ctx)

	// Verificar resultados
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, stats, result)

	// Verificar que se llamaron los métodos esperados
	mockRepo.AssertExpectations(t)
}
