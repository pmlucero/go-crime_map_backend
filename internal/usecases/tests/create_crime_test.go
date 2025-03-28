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
	"github.com/stretchr/testify/mock"
)

// MockCrimeRepository es un mock del repositorio para pruebas
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

func (m *MockCrimeRepository) List(ctx context.Context, page, limit int, startDate, endDate *time.Time, crimeType, status *string) ([]entities.Crime, int64, error) {
	args := m.Called(ctx, page, limit, startDate, endDate, crimeType, status)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entities.Crime), args.Get(1).(int64), args.Error(2)
}

func (m *MockCrimeRepository) GetStats(ctx context.Context) (*entities.CrimeStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CrimeStats), args.Error(1)
}

func TestCreateCrimeUseCase_Execute(t *testing.T) {
	mockRepo := new(mocks.MockCrimeRepository)
	useCase := usecases.NewCreateCrimeUseCase(mockRepo)

	tests := []struct {
		name          string
		input         domain_usecases.CreateCrimeInput
		expectedError string
		setupMock     func()
	}{
		{
			name: "creación exitosa de delito",
			input: domain_usecases.CreateCrimeInput{
				Title:       "Robo a mano armada",
				Type:        "ROBO",
				Description: "Robo a mano armada en comercio",
				Latitude:    -34.603722,
				Longitude:   -58.381592,
				Address:     "Av. Corrientes 1234",
			},
			setupMock: func() {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Crime")).Return(nil)
			},
		},
		{
			name: "error - tipo de delito vacío",
			input: domain_usecases.CreateCrimeInput{
				Title:       "Robo a mano armada",
				Description: "Robo a mano armada en comercio",
				Latitude:    -34.603722,
				Longitude:   -58.381592,
				Address:     "Av. Corrientes 1234",
			},
			expectedError: "el tipo es requerido",
		},
		{
			name: "error - descripción vacía",
			input: domain_usecases.CreateCrimeInput{
				Title:     "Robo a mano armada",
				Type:      "ROBO",
				Latitude:  -34.603722,
				Longitude: -58.381592,
				Address:   "Av. Corrientes 1234",
			},
			expectedError: "la descripción es requerida",
		},
		{
			name: "error - latitud inválida",
			input: domain_usecases.CreateCrimeInput{
				Title:       "Robo a mano armada",
				Type:        "ROBO",
				Description: "Robo a mano armada en comercio",
				Latitude:    91.0,
				Longitude:   -58.381592,
				Address:     "Av. Corrientes 1234",
			},
			expectedError: "la latitud debe estar entre -90 y 90",
		},
		{
			name: "error - longitud inválida",
			input: domain_usecases.CreateCrimeInput{
				Title:       "Robo a mano armada",
				Type:        "ROBO",
				Description: "Robo a mano armada en comercio",
				Latitude:    -34.603722,
				Longitude:   181.0,
				Address:     "Av. Corrientes 1234",
			},
			expectedError: "la longitud debe estar entre -180 y 180",
		},
		{
			name: "error - fallo en el repositorio",
			input: domain_usecases.CreateCrimeInput{
				Title:       "Robo a mano armada",
				Type:        "ROBO",
				Description: "Robo a mano armada en comercio",
				Latitude:    -34.603722,
				Longitude:   -58.381592,
				Address:     "Av. Corrientes 1234",
			},
			expectedError: "error al crear el delito: assert.AnError general error for testing",
			setupMock: func() {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Crime")).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil

			if tt.setupMock != nil {
				tt.setupMock()
			}

			result, err := useCase.Execute(context.Background(), tt.input)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.ID)
			assert.Equal(t, tt.input.Title, result.Title)
			assert.Equal(t, tt.input.Type, result.Type)
			assert.Equal(t, tt.input.Description, result.Description)
			assert.Equal(t, tt.input.Latitude, result.Location.Latitude)
			assert.Equal(t, tt.input.Longitude, result.Location.Longitude)
			assert.Equal(t, tt.input.Address, result.Location.Address)
			assert.Equal(t, string(entities.CrimeStatusActive), result.Status)

			mockRepo.AssertExpectations(t)
		})
	}
}
