package usecases

import (
	"context"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/usecases"
	"go-crime_map_backend/internal/mocks"
	"go-crime_map_backend/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestListCrimesUseCase_Execute(t *testing.T) {
	// Crear mock del repositorio
	mockRepo := new(mocks.MockCrimeRepository)
	useCase := NewListCrimesUseCase(mockRepo)

	// Crear contexto
	ctx := context.Background()

	// Crear datos de prueba
	crimes := []entities.Crime{
		{
			ID:          "1",
			Title:       "Robo de auto",
			Description: "Me robaron el auto estacionado",
			Type:        "ROBO",
			Status:      "ACTIVE",
			Location: entities.Location{
				Latitude:      -34.603722,
				Longitude:     -58.381592,
				Address:       "Av. 9 de Julio",
				AddressNumber: utils.StringPtr("1000"),
				City:          utils.StringPtr("Buenos Aires"),
				Province:      utils.StringPtr("CABA"),
				Country:       utils.StringPtr("Argentina"),
				ZipCode:       utils.StringPtr("C1043"),
			},
		},
		{
			ID:          "2",
			Title:       "Robo de moto",
			Description: "Me robaron la moto estacionada",
			Type:        "ROBO",
			Status:      "ACTIVE",
			Location: entities.Location{
				Latitude:      -34.603722,
				Longitude:     -58.381592,
				Address:       "Av. Corrientes",
				AddressNumber: utils.StringPtr("2000"),
				City:          utils.StringPtr("Buenos Aires"),
				Province:      utils.StringPtr("CABA"),
				Country:       utils.StringPtr("Argentina"),
				ZipCode:       utils.StringPtr("C1043"),
			},
		},
	}

	tests := []struct {
		name          string
		input         usecases.ListCrimesParams
		mockSetup     func()
		expectedError error
		expectedCount int
	}{
		{
			name: "Listar todos los delitos",
			input: usecases.ListCrimesParams{
				Page:  1,
				Limit: 10,
			},
			mockSetup: func() {
				mockRepo.On("List", ctx, 1, 10).Return(crimes, int64(2), nil)
			},
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name: "Listar delitos con límite personalizado",
			input: usecases.ListCrimesParams{
				Page:  1,
				Limit: 5,
			},
			mockSetup: func() {
				mockRepo.On("List", ctx, 1, 5).Return(crimes[:1], int64(1), nil)
			},
			expectedError: nil,
			expectedCount: 1,
		},
		{
			name: "Listar delitos con página personalizada",
			input: usecases.ListCrimesParams{
				Page:  2,
				Limit: 1,
			},
			mockSetup: func() {
				mockRepo.On("List", ctx, 2, 1).Return(crimes[1:], int64(1), nil)
			},
			expectedError: nil,
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar el mock
			tt.mockSetup()

			// Ejecutar el caso de uso
			result, err := useCase.Execute(ctx, tt.input)

			// Verificar el resultado
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedCount, len(result.Items))
			assert.Equal(t, int64(tt.expectedCount), result.Total)

			// Verificar que se llamó al mock
			mockRepo.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
