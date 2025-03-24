package usecases

import (
	"context"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"

	"github.com/stretchr/testify/assert"
)

func TestGetCrimeStatsUseCase_Execute(t *testing.T) {
	// Crear mock del repositorio
	mockRepo := new(MockCrimeRepository)
	useCase := NewGetCrimeStatsUseCase(mockRepo)

	// Crear contexto
	ctx := context.Background()

	// Crear datos de prueba
	crimes := []*entities.Crime{
		{
			ID:          "1",
			Type:        "ROBO",
			Description: "Robo a mano armada",
			Location: entities.Location{
				Latitude:  -34.603722,
				Longitude: -58.381592,
			},
			Status:    entities.CrimeStatusActive,
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:          "2",
			Type:        "ASALTO",
			Description: "Asalto a comercio",
			Location: entities.Location{
				Latitude:  -34.608722,
				Longitude: -58.382592,
			},
			Status:    entities.CrimeStatusActive,
			CreatedAt: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:          "3",
			Type:        "ROBO",
			Description: "Robo de vehículo",
			Location: entities.Location{
				Latitude:  -34.613722,
				Longitude: -58.383592,
			},
			Status:    entities.CrimeStatusInactive,
			CreatedAt: time.Now(),
		},
	}

	tests := []struct {
		name          string
		input         GetCrimeStatsInput
		mockSetup     func()
		expectedError error
		validateStats func(*testing.T, *CrimeStats)
	}{
		{
			name: "Obtener estadísticas completas",
			input: GetCrimeStatsInput{
				Limit: 5,
			},
			mockSetup: func() {
				mockRepo.On("List", ctx, repositories.ListCrimesFilter{
					Limit:  5,
					Offset: 0,
				}).Return(crimes, nil)
			},
			expectedError: nil,
			validateStats: func(t *testing.T, stats *CrimeStats) {
				assert.Equal(t, 3, stats.TotalCrimes)
				assert.Equal(t, 2, stats.CrimesByType["ROBO"])
				assert.Equal(t, 1, stats.CrimesByType["ASALTO"])
				assert.Equal(t, 2, stats.CrimesByStatus[entities.CrimeStatusActive])
				assert.Equal(t, 1, stats.CrimesByStatus[entities.CrimeStatusInactive])
				assert.Len(t, stats.MostCommonTypes, 2)
				assert.Equal(t, "ROBO", stats.MostCommonTypes[0])
				assert.Equal(t, "ASALTO", stats.MostCommonTypes[1])
			},
		},
		{
			name: "Error en rango de fechas inválido",
			input: GetCrimeStatsInput{
				StartDate: time.Now(),
				EndDate:   time.Now().Add(-24 * time.Hour),
				Limit:     5,
			},
			mockSetup:     func() {},
			expectedError: ErrInvalidDateRange,
			validateStats: nil,
		},
		{
			name: "Límite inválido se ajusta al valor por defecto",
			input: GetCrimeStatsInput{
				Limit: 0,
			},
			mockSetup: func() {
				mockRepo.On("List", ctx, repositories.ListCrimesFilter{
					Limit:  5,
					Offset: 0,
				}).Return(crimes, nil)
			},
			expectedError: nil,
			validateStats: func(t *testing.T, stats *CrimeStats) {
				assert.Equal(t, 3, stats.TotalCrimes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mock
			tt.mockSetup()

			// Ejecutar caso de uso
			stats, err := useCase.Execute(ctx, tt.input)

			// Verificar resultados
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, stats)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, stats)

			// Validar estadísticas
			if tt.validateStats != nil {
				tt.validateStats(t, stats)
			}

			// Verificar que se llamó al mock
			mockRepo.AssertExpectations(t)
		})
	}
}
