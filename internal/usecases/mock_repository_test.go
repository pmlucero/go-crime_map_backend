package usecases

import (
	"context"
	"time"

	"go-crime_map_backend/internal/domain/entities"

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
