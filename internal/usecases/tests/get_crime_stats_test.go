package tests

import (
	"context"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories/mocks"
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

	t.Run("Estadísticas correctas", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewGetCrimeStatsUseCase(mockRepo)
		mockRepo.On("GetStats", ctx).Return(stats, nil)
		// Act
		stats, err := useCase.Execute(ctx)
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, stats, stats)
	})

	t.Run("Error del repositorio", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewGetCrimeStatsUseCase(mockRepo)
		expectedErr := assert.AnError
		mockRepo.On("GetStats", ctx).Return(nil, expectedErr)
		// Act
		stats, err := useCase.Execute(ctx)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Equal(t, expectedErr, err)
	})
}
