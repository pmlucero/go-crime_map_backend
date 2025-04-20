package tests

import (
	"context"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/mocks"
	"go-crime_map_backend/internal/usecases"

	"github.com/stretchr/testify/assert"
)

func TestGetMostCommonTypes(t *testing.T) {
	tests := []struct {
		name         string
		crimesByType map[string]int
		expected     []string
	}{
		{
			name: "obtener top 5 tipos más comunes",
			crimesByType: map[string]int{
				"ROBO":      50,
				"HURTO":     30,
				"ASALTO":    20,
				"VIOLENCIA": 15,
				"FRAUDE":    10,
				"OTRO":      5,
			},
			expected: []string{"ROBO", "HURTO", "ASALTO", "VIOLENCIA", "FRAUDE"},
		},
		{
			name: "menos de 5 tipos",
			crimesByType: map[string]int{
				"ROBO":   50,
				"HURTO":  30,
				"ASALTO": 20,
			},
			expected: []string{"ROBO", "HURTO", "ASALTO"},
		},
		{
			name:         "mapa vacío",
			crimesByType: map[string]int{},
			expected:     []string{},
		},
		{
			name: "empates en cantidad",
			crimesByType: map[string]int{
				"ROBO":   50,
				"HURTO":  50,
				"ASALTO": 30,
			},
			expected: []string{"HURTO", "ROBO", "ASALTO"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := usecases.GetMostCommonTypes(tt.crimesByType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetCrimeStatsUseCase_Execute(t *testing.T) {
	// Crear mock del repositorio
	mockRepo := new(mocks.MockCrimeRepository)
	useCase := usecases.NewGetCrimeStatsUseCase(mockRepo)

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
		CrimesByAddress: map[string]int64{},
		LastUpdate:      time.Now(),
	}

	tests := []struct {
		name          string
		mockSetup     func()
		expectedStats *entities.CrimeStats
		expectedError error
	}{
		{
			name: "obtener estadísticas exitosamente",
			mockSetup: func() {
				mockRepo.On("GetStats", ctx).Return(stats, nil).Once()
			},
			expectedStats: stats,
			expectedError: nil,
		},
		{
			name: "error al obtener estadísticas",
			mockSetup: func() {
				mockRepo.On("GetStats", ctx).Return(nil, assert.AnError).Once()
			},
			expectedStats: nil,
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Limpiar mock y configurar para este test
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
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
			assert.Equal(t, tt.expectedStats.TotalCrimes, result.TotalCrimes)
			assert.Equal(t, tt.expectedStats.ActiveCrimes, result.ActiveCrimes)
			assert.Equal(t, tt.expectedStats.InactiveCrimes, result.InactiveCrimes)
			assert.Equal(t, tt.expectedStats.CrimesByType, result.CrimesByType)
			assert.Equal(t, tt.expectedStats.CrimesByStatus, result.CrimesByStatus)
			assert.Equal(t, tt.expectedStats.CrimesByLocation, result.CrimesByLocation)
			assert.Equal(t, tt.expectedStats.CrimesByAddress, result.CrimesByAddress)
		})
	}
}
