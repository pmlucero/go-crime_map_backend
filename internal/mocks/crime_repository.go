package mocks

import (
	"context"
	"time"

	"go-crime_map_backend/internal/domain/entities"

	"github.com/stretchr/testify/mock"
)

// MockCrimeRepository es un mock del repositorio para pruebas
type MockCrimeRepository struct {
	*mock.Mock
}

// NewMockCrimeRepository crea una nueva instancia del mock
func NewMockCrimeRepository() *MockCrimeRepository {
	return &MockCrimeRepository{
		Mock: &mock.Mock{},
	}
}

// Create implementa la interfaz CrimeRepository
func (m *MockCrimeRepository) Create(ctx context.Context, crime *entities.Crime) error {
	args := m.Mock.Called(ctx, crime)
	return args.Error(0)
}

// GetByID implementa la interfaz CrimeRepository
func (m *MockCrimeRepository) GetByID(ctx context.Context, id string) (*entities.Crime, error) {
	args := m.Mock.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Crime), args.Error(1)
}

// GetAll implementa la interfaz CrimeRepository
func (m *MockCrimeRepository) GetAll(ctx context.Context) ([]*entities.Crime, error) {
	args := m.Mock.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Crime), args.Error(1)
}

// Update implementa la interfaz CrimeRepository
func (m *MockCrimeRepository) Update(ctx context.Context, crime *entities.Crime) error {
	args := m.Mock.Called(ctx, crime)
	return args.Error(0)
}

// Delete implementa la interfaz CrimeRepository
func (m *MockCrimeRepository) Delete(ctx context.Context, id string) error {
	args := m.Mock.Called(ctx, id)
	return args.Error(0)
}

// List implementa la interfaz CrimeRepository
func (m *MockCrimeRepository) List(ctx context.Context, page, limit int, startDate, endDate *time.Time, crimeType, status *string) ([]entities.Crime, int64, error) {
	args := m.Mock.Called(ctx, page, limit, startDate, endDate, crimeType, status)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entities.Crime), args.Get(1).(int64), args.Error(2)
}

// GetStats implementa la interfaz CrimeRepository
func (m *MockCrimeRepository) GetStats(ctx context.Context) (*entities.CrimeStats, error) {
	args := m.Mock.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CrimeStats), args.Error(1)
}
