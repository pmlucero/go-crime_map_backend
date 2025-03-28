package usecases

import (
	"context"
	"testing"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/mocks"

	"github.com/stretchr/testify/assert"
)

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
			"HURTO":  30,
			"ASALTO": 20,
		},
		CrimesByStatus: map[string]int64{
			"ACTIVE":   80,
			"INACTIVE": 20,
		},
		CrimesByLocation: map[string]int64{
			"CENTRO": 40,
			"NORTE":  30,
			"SUR":    30,
		},
	}

	tests := []struct {
		name          string
		mockSetup     func()
		expectedError error
	}{
		{
			name: "obtener estadísticas exitosamente",
			mockSetup: func() {
				mockRepo.On("GetStats", ctx).Return(stats, nil)
			},
			expectedError: nil,
		},
		{
			name: "error al obtener estadísticas",
			mockSetup: func() {
				mockRepo.On("GetStats", ctx).Return(nil, assert.AnError)
			},
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar el mock
			tt.mockSetup()

			// Ejecutar el caso de uso
			result, err := useCase.Execute(ctx)

			// Verificar el resultado
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, stats.TotalCrimes, result.TotalCrimes)
			assert.Equal(t, stats.ActiveCrimes, result.ActiveCrimes)
			assert.Equal(t, stats.InactiveCrimes, result.InactiveCrimes)
			assert.Equal(t, stats.CrimesByType, result.CrimesByType)
			assert.Equal(t, stats.CrimesByStatus, result.CrimesByStatus)
			assert.Equal(t, stats.CrimesByLocation, result.CrimesByLocation)
		})
	}
}
