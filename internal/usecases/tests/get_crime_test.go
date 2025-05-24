package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories/mocks"
	"go-crime_map_backend/internal/usecases"

	"github.com/stretchr/testify/assert"
)

func TestGetCrimeUseCase_Execute(t *testing.T) {
	mockRepo := mocks.NewCrimeRepository(t)
	useCase := usecases.NewGetCrimeUseCase(mockRepo)

	ctx := context.Background()

	// Crear un delito de ejemplo para los tests
	crime := &entities.Crime{
		ID:          1,
		UUID:        "1",
		Title:       "Robo a mano armada",
		Type:        "ROBO",
		Description: "Robo a mano armada en la vía pública",
		Location: entities.Location{
			Latitude:      -34.603722,
			Longitude:     -58.381592,
			Address:       "Av. Corrientes 1234",
			AddressNumber: "1234",
			City:          "Buenos Aires",
			Province:      "CABA",
			Country:       "Argentina",
			ZipCode:       "C1043AAZ",
		},
		Status:    string(entities.CrimeStatusActive),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name          string
		id            string
		mockSetup     func()
		expectedError string
		expectedCrime *entities.Crime
	}{
		{
			name:          "Error - ID vacío",
			id:            "",
			mockSetup:     func() {},
			expectedError: "el UUID del delito es requerido",
		},
		{
			name: "Error - error al obtener el delito",
			id:   "1",
			mockSetup: func() {
				mockRepo.On("GetByUUID", ctx, "1").Return(nil, errors.New("error de base de datos"))
			},
			expectedError: "error al obtener el delito: error de base de datos",
		},
		{
			name: "Error - delito no encontrado",
			id:   "1",
			mockSetup: func() {
				mockRepo.On("GetByUUID", ctx, "1").Return(nil, nil)
			},
			expectedError: "delito no encontrado",
		},
		{
			name: "Obtener delito exitosamente",
			id:   "1",
			mockSetup: func() {
				mockRepo.On("GetByUUID", ctx, "1").Return(crime, nil)
			},
			expectedCrime: crime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.mockSetup()
			// Act
			result, err := useCase.Execute(ctx, tt.id)
			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, result)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCrime, result)
			mockRepo.AssertExpectations(t)
		})
	}
}
