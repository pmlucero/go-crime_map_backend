package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/mocks"
	"go-crime_map_backend/internal/usecases"

	"github.com/stretchr/testify/assert"
)

func TestDeleteCrimeUseCase_Execute(t *testing.T) {
	mockRepo := new(mocks.MockCrimeRepository)
	useCase := usecases.NewDeleteCrimeUseCase(mockRepo)

	ctx := context.Background()

	activeCrime := &entities.Crime{
		ID:          "1",
		Type:        "ROBO",
		Description: "Robo a mano armada",
		Location: entities.Location{
			Latitude:  -34.603722,
			Longitude: -58.381592,
		},
		Status:    string(entities.CrimeStatusActive),
		CreatedAt: time.Now(),
	}

	deletedCrime := &entities.Crime{
		ID:          "2",
		Type:        "ASALTO",
		Description: "Asalto a comercio",
		Location: entities.Location{
			Latitude:  -34.608722,
			Longitude: -58.382592,
		},
		Status:    string(entities.CrimeStatusDeleted),
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name          string
		id            string
		mockSetup     func()
		expectedError string
	}{
		{
			name:          "Error - ID vacío",
			id:            "",
			mockSetup:     func() {},
			expectedError: "el ID es requerido",
		},
		{
			name: "Error - error al obtener el delito",
			id:   "1",
			mockSetup: func() {
				mockRepo.On("GetByID", ctx, "1").Return(nil, errors.New("error de base de datos"))
			},
			expectedError: "error al obtener el delito: error de base de datos",
		},
		{
			name: "Error - delito no encontrado",
			id:   "1",
			mockSetup: func() {
				mockRepo.On("GetByID", ctx, "1").Return(nil, nil)
			},
			expectedError: "delito no encontrado",
		},
		{
			name: "Error - delito ya eliminado",
			id:   "2",
			mockSetup: func() {
				mockRepo.On("GetByID", ctx, "2").Return(deletedCrime, nil)
			},
			expectedError: "el delito ya fue eliminado",
		},
		{
			name: "Error - error al eliminar el delito",
			id:   "1",
			mockSetup: func() {
				mockRepo.On("GetByID", ctx, "1").Return(activeCrime, nil)
				mockRepo.On("Delete", ctx, "1").Return(errors.New("error al eliminar"))
			},
			expectedError: "error al eliminar el delito: error al eliminar",
		},
		{
			name: "Eliminar delito exitosamente",
			id:   "1",
			mockSetup: func() {
				mockRepo.On("GetByID", ctx, "1").Return(activeCrime, nil)
				mockRepo.On("Delete", ctx, "1").Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Limpiar mock y configurar para este test
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.mockSetup()

			err := useCase.Execute(ctx, tt.id)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				return
			}

			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
