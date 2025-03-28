package usecases

import (
	"context"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
	"go-crime_map_backend/internal/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateCrimeStatusUseCase_Execute(t *testing.T) {
	// Crear mock del repositorio
	mockRepo := new(mocks.MockCrimeRepository)
	useCase := NewUpdateCrimeStatusUseCase(mockRepo)

	// Crear contexto
	ctx := context.Background()

	// Crear datos de prueba
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

	inactiveCrime := &entities.Crime{
		ID:          "2",
		Type:        "ASALTO",
		Description: "Asalto a comercio",
		Location: entities.Location{
			Latitude:  -34.608722,
			Longitude: -58.382592,
		},
		Status:    string(entities.CrimeStatusInactive),
		CreatedAt: time.Now(),
	}

	deletedCrime := &entities.Crime{
		ID:          "3",
		Type:        "ROBO",
		Description: "Robo de vehículo",
		Location: entities.Location{
			Latitude:  -34.613722,
			Longitude: -58.383592,
		},
		Status:    string(entities.CrimeStatusDeleted),
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name          string
		input         UpdateCrimeStatusInput
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Actualizar estado de activo a inactivo",
			input: UpdateCrimeStatusInput{
				ID:        "1",
				NewStatus: string(entities.CrimeStatusInactive),
			},
			mockSetup: func() {
				mockRepo.On("GetByID", ctx, "1").Return(activeCrime, nil)
				mockRepo.On("Update", ctx, mock.MatchedBy(func(crime *entities.Crime) bool {
					return crime.ID == "1" && crime.Status == string(entities.CrimeStatusInactive)
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Actualizar estado de inactivo a activo",
			input: UpdateCrimeStatusInput{
				ID:        "2",
				NewStatus: string(entities.CrimeStatusActive),
			},
			mockSetup: func() {
				mockRepo.On("GetByID", ctx, "2").Return(inactiveCrime, nil)
				mockRepo.On("Update", ctx, mock.MatchedBy(func(crime *entities.Crime) bool {
					return crime.ID == "2" && crime.Status == string(entities.CrimeStatusActive)
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Error al intentar actualizar un delito eliminado",
			input: UpdateCrimeStatusInput{
				ID:        "3",
				NewStatus: string(entities.CrimeStatusActive),
			},
			mockSetup: func() {
				mockRepo.On("GetByID", ctx, "3").Return(deletedCrime, nil)
			},
			expectedError: ErrCrimeAlreadyDeleted,
		},
		{
			name: "Error al intentar cambiar a estado eliminado",
			input: UpdateCrimeStatusInput{
				ID:        "1",
				NewStatus: string(entities.CrimeStatusDeleted),
			},
			mockSetup: func() {
				mockRepo.On("GetByID", ctx, "1").Return(activeCrime, nil)
			},
			expectedError: ErrInvalidStatusTransition,
		},
		{
			name: "Error al intentar cambiar a estado inválido",
			input: UpdateCrimeStatusInput{
				ID:        "1",
				NewStatus: "INVALID_STATUS",
			},
			mockSetup: func() {
				mockRepo.On("GetByID", ctx, "1").Return(activeCrime, nil)
			},
			expectedError: ErrInvalidStatusTransition,
		},
		{
			name: "Error al no encontrar el delito",
			input: UpdateCrimeStatusInput{
				ID:        "4",
				NewStatus: string(entities.CrimeStatusInactive),
			},
			mockSetup: func() {
				mockRepo.On("GetByID", ctx, "4").Return(nil, repositories.ErrNotFound)
			},
			expectedError: repositories.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mock
			tt.mockSetup()

			// Ejecutar caso de uso
			err := useCase.Execute(ctx, tt.input)

			// Verificar resultados
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				return
			}

			assert.NoError(t, err)

			// Verificar que se llamó al mock
			mockRepo.AssertExpectations(t)
		})
	}
}
