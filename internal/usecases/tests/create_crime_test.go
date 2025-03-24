package tests

import (
	"context"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
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

func TestCreateCrimeUseCase_Execute(t *testing.T) {
	mockRepo := new(MockCrimeRepository)
	useCase := usecases.NewCreateCrimeUseCase(mockRepo)

	// Datos de prueba comunes
	validLocation := usecases.Location{
		Latitude:  -34.603722,
		Longitude: -58.381592,
		Address:   "Av. Corrientes 1234",
	}

	tests := []struct {
		name          string
		input         usecases.CreateCrimeInput
		expectedError string
		setupMock     func()
	}{
		{
			name: "creación exitosa de delito",
			input: usecases.CreateCrimeInput{
				Type:        "ROBO",
				Description: "Robo a mano armada",
				Location:    validLocation,
				Date:        time.Now().Add(-1 * time.Hour),
			},
			setupMock: func() {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Crime")).Return(nil)
				mockRepo.On("GetAll", mock.Anything).Return([]*entities.Crime{}, nil)
			},
		},
		{
			name: "error - tipo de delito vacío",
			input: usecases.CreateCrimeInput{
				Type:        "",
				Description: "Robo a mano armada",
				Location:    validLocation,
				Date:        time.Now().Add(-1 * time.Hour),
			},
			expectedError: "el tipo de delito es inválido",
		},
		{
			name: "error - descripción vacía",
			input: usecases.CreateCrimeInput{
				Type:        "ROBO",
				Description: "",
				Location:    validLocation,
				Date:        time.Now().Add(-1 * time.Hour),
			},
			expectedError: "la descripción es requerida",
		},
		{
			name: "error - fecha futura",
			input: usecases.CreateCrimeInput{
				Type:        "ROBO",
				Description: "Robo a mano armada",
				Location:    validLocation,
				Date:        time.Now().Add(24 * time.Hour),
			},
			expectedError: "la fecha del delito no puede ser futura",
		},
		{
			name: "error - latitud inválida",
			input: usecases.CreateCrimeInput{
				Type:        "ROBO",
				Description: "Robo a mano armada",
				Location: usecases.Location{
					Latitude:  91.0, // Latitud inválida
					Longitude: -58.381592,
					Address:   "Av. Corrientes 1234",
				},
				Date: time.Now().Add(-1 * time.Hour),
			},
			expectedError: "latitud inválida",
		},
		{
			name: "error - longitud inválida",
			input: usecases.CreateCrimeInput{
				Type:        "ROBO",
				Description: "Robo a mano armada",
				Location: usecases.Location{
					Latitude:  -34.603722,
					Longitude: 181.0, // Longitud inválida
					Address:   "Av. Corrientes 1234",
				},
				Date: time.Now().Add(-1 * time.Hour),
			},
			expectedError: "longitud inválida",
		},
		{
			name: "error - fallo en el repositorio",
			input: usecases.CreateCrimeInput{
				Type:        "ROBO",
				Description: "Robo a mano armada",
				Location:    validLocation,
				Date:        time.Now().Add(-1 * time.Hour),
			},
			expectedError: "assert.AnError general error for testing",
			setupMock: func() {
				mockRepo.On("GetAll", mock.Anything).Return([]*entities.Crime{}, nil)
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Crime")).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Limpiar las expectativas del mock
			mockRepo.ExpectedCalls = nil

			// Configurar el mock si es necesario
			if tt.setupMock != nil {
				tt.setupMock()
			} else if tt.expectedError == "" {
				// Si no hay error esperado y no hay setup específico, configurar el mock por defecto
				mockRepo.On("GetAll", mock.Anything).Return([]*entities.Crime{}, nil)
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Crime")).Return(nil)
			}

			// Ejecutar el caso de uso
			result, err := useCase.Execute(context.Background(), tt.input)

			// Verificar el resultado
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.ID)
			assert.Equal(t, tt.input.Type, result.Type)
			assert.Equal(t, tt.input.Description, result.Description)
			assert.Equal(t, tt.input.Location.Latitude, result.Location.Latitude)
			assert.Equal(t, tt.input.Location.Longitude, result.Location.Longitude)
			assert.Equal(t, tt.input.Location.Address, result.Location.Address)
			assert.Equal(t, tt.input.Date, result.Date)

			// Verificar que se llamó al repositorio
			mockRepo.AssertExpectations(t)
		})
	}
}
