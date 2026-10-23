package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
	"go-crime_map_backend/internal/domain/repositories/mocks"
	domain_usecases "go-crime_map_backend/internal/domain/usecases"
	"go-crime_map_backend/internal/usecases"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateCrimeStatusUseCase_Execute(t *testing.T) {
	mockRepo := mocks.NewCrimeRepository(t)
	useCase := usecases.NewUpdateCrimeStatusUseCase(mockRepo)

	ctx := context.Background()

	activeCrime := &entities.Crime{
		ID:          1,
		UUID:        "1",
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
		ID:          2,
		UUID:        "2",
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
		ID:          3,
		UUID:        "3",
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
		input         domain_usecases.UpdateCrimeStatusInput
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Error - UUID vacío",
			input: domain_usecases.UpdateCrimeStatusInput{
				UUID:   "",
				Status: string(entities.CrimeStatusActive),
			},
			mockSetup:     func() {},
			expectedError: fmt.Errorf("el UUID es requerido"),
		},
		{
			name: "Error - Status vacío",
			input: domain_usecases.UpdateCrimeStatusInput{
				UUID:   "1",
				Status: "",
			},
			mockSetup:     func() {},
			expectedError: fmt.Errorf("el estado es requerido"),
		},
		{
			name: "Error genérico al obtener el delito",
			input: domain_usecases.UpdateCrimeStatusInput{
				UUID:   "1",
				Status: string(entities.CrimeStatusActive),
			},
			mockSetup: func() {
				mockRepo.On("GetByUUID", ctx, "1").Return(nil, errors.New("error de base de datos"))
			},
			expectedError: fmt.Errorf("error al obtener el delito: %w", errors.New("error de base de datos")),
		},
		{
			name: "Error - GetByUUID devuelve nil sin error",
			input: domain_usecases.UpdateCrimeStatusInput{
				UUID:   "1",
				Status: string(entities.CrimeStatusActive),
			},
			mockSetup: func() {
				mockRepo.On("GetByUUID", ctx, "1").Return(nil, nil)
			},
			expectedError: repositories.ErrNotFound,
		},
		{
			name: "Error al actualizar el delito",
			input: domain_usecases.UpdateCrimeStatusInput{
				UUID:   "1",
				Status: string(entities.CrimeStatusInactive),
			},
			mockSetup: func() {
				mockRepo.On("GetByUUID", ctx, "1").Return(activeCrime, nil)
				mockRepo.On("Update", ctx, mock.MatchedBy(func(crime *entities.Crime) bool {
					return crime.UUID == "1" && crime.Status == string(entities.CrimeStatusInactive)
				})).Return(errors.New("error al actualizar"))
			},
			expectedError: fmt.Errorf("error al actualizar el delito: %w", errors.New("error al actualizar")),
		},
		{
			name: "Actualizar estado de activo a inactivo",
			input: domain_usecases.UpdateCrimeStatusInput{
				UUID:   "1",
				Status: string(entities.CrimeStatusInactive),
			},
			mockSetup: func() {
				mockRepo.On("GetByUUID", ctx, "1").Return(activeCrime, nil)
				mockRepo.On("Update", ctx, mock.MatchedBy(func(crime *entities.Crime) bool {
					return crime.UUID == "1" && crime.Status == string(entities.CrimeStatusInactive)
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Actualizar estado de inactivo a activo",
			input: domain_usecases.UpdateCrimeStatusInput{
				UUID:   "2",
				Status: string(entities.CrimeStatusActive),
			},
			mockSetup: func() {
				mockRepo.On("GetByUUID", ctx, "2").Return(inactiveCrime, nil)
				mockRepo.On("Update", ctx, mock.MatchedBy(func(crime *entities.Crime) bool {
					return crime.UUID == "2" && crime.Status == string(entities.CrimeStatusActive)
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Error al intentar actualizar un delito eliminado",
			input: domain_usecases.UpdateCrimeStatusInput{
				UUID:   "3",
				Status: string(entities.CrimeStatusActive),
			},
			mockSetup: func() {
				mockRepo.On("GetByUUID", ctx, "3").Return(deletedCrime, nil)
			},
			expectedError: usecases.ErrCrimeAlreadyDeleted,
		},
		{
			name: "Error al intentar cambiar a estado eliminado",
			input: domain_usecases.UpdateCrimeStatusInput{
				UUID:   "1",
				Status: string(entities.CrimeStatusDeleted),
			},
			mockSetup: func() {
				mockRepo.On("GetByUUID", ctx, "1").Return(activeCrime, nil)
			},
			expectedError: usecases.ErrInvalidStatusTransition,
		},
		{
			name: "Error al intentar cambiar a estado inválido",
			input: domain_usecases.UpdateCrimeStatusInput{
				UUID:   "1",
				Status: "INVALID_STATUS",
			},
			mockSetup: func() {
				mockRepo.On("GetByUUID", ctx, "1").Return(activeCrime, nil)
			},
			expectedError: usecases.ErrInvalidStatusTransition,
		},
		{
			name: "Error al no encontrar el delito",
			input: domain_usecases.UpdateCrimeStatusInput{
				UUID:   "4",
				Status: string(entities.CrimeStatusInactive),
			},
			mockSetup: func() {
				mockRepo.On("GetByUUID", ctx, "4").Return(nil, repositories.ErrNotFound)
			},
			expectedError: repositories.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.mockSetup()
			// Act
			err := useCase.Execute(ctx, tt.input)
			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				return
			}
			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
