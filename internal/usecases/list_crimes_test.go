package usecases

import (
	"context"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCrimeRepository es un mock del repositorio de delitos
type MockCrimeRepository struct {
	mock.Mock
}

func (m *MockCrimeRepository) Create(ctx context.Context, crime *entities.Crime) error {
	args := m.Called(ctx, crime)
	return args.Error(0)
}

func (m *MockCrimeRepository) GetByID(ctx context.Context, id string) (*entities.Crime, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Crime), args.Error(1)
}

func (m *MockCrimeRepository) GetAll(ctx context.Context) ([]*entities.Crime, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Crime), args.Error(1)
}

func (m *MockCrimeRepository) Update(ctx context.Context, crime *entities.Crime) error {
	args := m.Called(ctx, crime)
	return args.Error(0)
}

func (m *MockCrimeRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCrimeRepository) List(ctx context.Context, filter repositories.ListCrimesFilter) ([]*entities.Crime, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Crime), args.Error(1)
}

func TestListCrimesUseCase_Execute(t *testing.T) {
	// Crear mock del repositorio
	mockRepo := new(MockCrimeRepository)
	useCase := NewListCrimesUseCase(mockRepo)

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
			Status:    "ACTIVO",
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
			Status:    "RESUELTO",
			CreatedAt: time.Now(),
		},
	}

	tests := []struct {
		name          string
		input         ListCrimesInput
		mockSetup     func()
		expectedError error
		expectedCount int
	}{
		{
			name: "Listar todos los delitos",
			input: ListCrimesInput{
				Limit:  10,
				Offset: 0,
			},
			mockSetup: func() {
				mockRepo.On("List", ctx, repositories.ListCrimesFilter{
					Limit:  10,
					Offset: 0,
				}).Return(crimes, nil)
			},
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name: "Listar delitos por tipo",
			input: ListCrimesInput{
				Type:   "ROBO",
				Limit:  10,
				Offset: 0,
			},
			mockSetup: func() {
				mockRepo.On("List", ctx, repositories.ListCrimesFilter{
					Type:   "ROBO",
					Limit:  10,
					Offset: 0,
				}).Return([]*entities.Crime{crimes[0]}, nil)
			},
			expectedError: nil,
			expectedCount: 1,
		},
		{
			name: "Listar delitos por rango de fechas",
			input: ListCrimesInput{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now(),
				Limit:     10,
				Offset:    0,
			},
			mockSetup: func() {
				mockRepo.On("List", ctx, repositories.ListCrimesFilter{
					StartDate: time.Now().Add(-24 * time.Hour),
					EndDate:   time.Now(),
					Limit:     10,
					Offset:    0,
				}).Return(crimes, nil)
			},
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name: "Listar delitos por ubicación",
			input: ListCrimesInput{
				Latitude:  -34.603722,
				Longitude: -58.381592,
				Radius:    1.0,
				Limit:     10,
				Offset:    0,
			},
			mockSetup: func() {
				mockRepo.On("List", ctx, repositories.ListCrimesFilter{
					Latitude:  -34.603722,
					Longitude: -58.381592,
					Radius:    1.0,
					Limit:     10,
					Offset:    0,
				}).Return(crimes, nil)
			},
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name: "Error en rango de fechas inválido",
			input: ListCrimesInput{
				StartDate: time.Now(),
				EndDate:   time.Now().Add(-24 * time.Hour),
				Limit:     10,
				Offset:    0,
			},
			mockSetup:     func() {},
			expectedError: ErrInvalidDateRange,
			expectedCount: 0,
		},
		{
			name: "Error en radio inválido",
			input: ListCrimesInput{
				Latitude:  -34.603722,
				Longitude: -58.381592,
				Radius:    0,
				Limit:     10,
				Offset:    0,
			},
			mockSetup:     func() {},
			expectedError: ErrInvalidRadius,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mock
			tt.mockSetup()

			// Ejecutar caso de uso
			result, err := useCase.Execute(ctx, tt.input)

			// Verificar resultados
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, tt.expectedCount)

			// Verificar que se llamó al mock
			mockRepo.AssertExpectations(t)
		})
	}
}
