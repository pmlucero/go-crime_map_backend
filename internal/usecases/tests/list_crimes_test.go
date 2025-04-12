package tests

import (
	"context"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	domain_usecases "go-crime_map_backend/internal/domain/usecases"
	"go-crime_map_backend/internal/mocks"
	"go-crime_map_backend/internal/usecases"

	"github.com/stretchr/testify/assert"
)

func TestListCrimesUseCase_Execute(t *testing.T) {
	// Crear mock del repositorio
	mockRepo := new(mocks.MockCrimeRepository)
	useCase := usecases.NewListCrimesUseCase(mockRepo)

	// Crear contexto
	ctx := context.Background()

	// Crear datos de prueba
	crimes := []entities.Crime{
		{
			ID:          "1",
			Description: "Robo a mano armada",
			Type:        "ROBO",
			Status:      "ACTIVE",
			Location: entities.Location{
				Latitude:  -34.603722,
				Longitude: -58.381592,
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          "2",
			Description: "Asalto a comercio",
			Type:        "ASALTO",
			Status:      "ACTIVE",
			Location: entities.Location{
				Latitude:  -34.608722,
				Longitude: -58.382592,
			},
			CreatedAt: time.Now(),
		},
	}

	tests := []struct {
		name          string
		params        domain_usecases.ListCrimesParams
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Listar todos los delitos",
			params: domain_usecases.ListCrimesParams{
				Page:  1,
				Limit: 10,
			},
			mockSetup: func() {
				mockRepo.On("List", ctx, 1, 10).Return(crimes, int64(2), nil)
			},
			expectedError: nil,
		},
		{
			name: "Error en rango de fechas inválido",
			params: domain_usecases.ListCrimesParams{
				Page:      1,
				Limit:     10,
				StartDate: time.Now().Add(24 * time.Hour),
				EndDate:   time.Now(),
			},
			mockSetup:     func() {},
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar el mock
			tt.mockSetup()

			// Ejecutar el caso de uso
			result, err := useCase.Execute(ctx, tt.params)

			// Verificar el error
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Items, len(crimes))
			}

			// Verificar que se llamó al mock
			mockRepo.AssertExpectations(t)
		})
	}
}
