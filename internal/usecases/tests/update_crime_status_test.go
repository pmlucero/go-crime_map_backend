package tests

import (
	"context"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
	domain_usecases "go-crime_map_backend/internal/domain/usecases"
	"go-crime_map_backend/internal/mocks"
	"go-crime_map_backend/internal/usecases"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateCrimeStatusUseCase_Execute(t *testing.T) {
	mockRepo := new(mocks.MockCrimeRepository)
	useCase := usecases.NewUpdateCrimeStatusUseCase(mockRepo)

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
		input         domain_usecases.UpdateCrimeStatusInput
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Actualizar estado de activo a inactivo",
			input: domain_usecases.UpdateCrimeStatusInput{
				ID:     "1",
				Status: string(entities.CrimeStatusInactive),
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
			input: domain_usecases.UpdateCrimeStatusInput{
				ID:     "2",
				Status: string(entities.CrimeStatusActive),
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
			input: domain_usecases.UpdateCrimeStatusInput{
				ID:     "3",
				Status: string(entities.CrimeStatusActive),
			},
			mockSetup: func() {
				mockRepo.On("GetByID", ctx, "3").Return(deletedCrime, nil)
			},
			expectedError: usecases.ErrCrimeAlreadyDeleted,
		},
		{
			name: "Error al intentar cambiar a estado eliminado",
			input: domain_usecases.UpdateCrimeStatusInput{
				ID:     "1",
				Status: string(entities.CrimeStatusDeleted),
			},
			mockSetup: func() {
				mockRepo.On("GetByID", ctx, "1").Return(activeCrime, nil)
			},
			expectedError: usecases.ErrInvalidStatusTransition,
		},
		{
			name: "Error al intentar cambiar a estado inválido",
			input: domain_usecases.UpdateCrimeStatusInput{
				ID:     "1",
				Status: "INVALID_STATUS",
			},
			mockSetup: func() {
				mockRepo.On("GetByID", ctx, "1").Return(activeCrime, nil)
			},
			expectedError: usecases.ErrInvalidStatusTransition,
		},
		{
			name: "Error al no encontrar el delito",
			input: domain_usecases.UpdateCrimeStatusInput{
				ID:     "4",
				Status: string(entities.CrimeStatusInactive),
			},
			mockSetup: func() {
				mockRepo.On("GetByID", ctx, "4").Return(nil, repositories.ErrNotFound)
			},
			expectedError: repositories.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := useCase.Execute(ctx, tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				return
			}

			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
