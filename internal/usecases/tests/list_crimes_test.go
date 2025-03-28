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
	mockRepo := new(mocks.MockCrimeRepository)
	useCase := usecases.NewListCrimesUseCase(mockRepo)

	ctx := context.Background()

	crimes := []entities.Crime{
		{
			ID:          "1",
			Type:        "ROBO",
			Description: "Robo a mano armada",
			Location: entities.Location{
				Latitude:  -34.603722,
				Longitude: -58.381592,
			},
			Status:    "ACTIVE",
			CreatedAt: time.Now(),
		},
		{
			ID:          "2",
			Type:        "ASALTO",
			Description: "Asalto a comercio",
			Location: entities.Location{
				Latitude:  -34.608722,
				Longitude: -58.382592,
			},
			Status:    "ACTIVE",
			CreatedAt: time.Now(),
		},
	}

	tests := []struct {
		name          string
		input         domain_usecases.ListCrimesParams
		mockSetup     func()
		expectedError error
		expectedCount int
	}{
		{
			name: "Listar todos los delitos",
			input: domain_usecases.ListCrimesParams{
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
			name: "Error en rango de fechas inválido",
			input: domain_usecases.ListCrimesParams{
				Page:  1,
				Limit: 10,
			},
			mockSetup:     func() {},
			expectedError: usecases.ErrInvalidDateRange,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := useCase.Execute(ctx, tt.input)

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

			mockRepo.AssertExpectations(t)
		})
	}
}
