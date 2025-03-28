package mocks

import (
	"context"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/usecases"

	"github.com/stretchr/testify/mock"
)

// MockCreateCrimeUseCase es un mock del caso de uso de creación de delitos
type MockCreateCrimeUseCase struct {
	mock.Mock
}

func (m *MockCreateCrimeUseCase) Execute(ctx context.Context, input usecases.CreateCrimeInput) (*entities.Crime, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Crime), args.Error(1)
}

// MockListCrimesUseCase es un mock del caso de uso de listado de delitos
type MockListCrimesUseCase struct {
	mock.Mock
}

func (m *MockListCrimesUseCase) Execute(ctx context.Context, params usecases.ListCrimesParams) (*entities.CrimeList, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CrimeList), args.Error(1)
}

// MockUpdateCrimeStatusUseCase es un mock del caso de uso de actualización de estado
type MockUpdateCrimeStatusUseCase struct {
	mock.Mock
}

func (m *MockUpdateCrimeStatusUseCase) Execute(ctx context.Context, input usecases.UpdateCrimeStatusInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

// MockDeleteCrimeUseCase es un mock del caso de uso de eliminación de delitos
type MockDeleteCrimeUseCase struct {
	mock.Mock
}

func (m *MockDeleteCrimeUseCase) Execute(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockGetCrimeStatsUseCase es un mock del caso de uso de obtención de estadísticas
type MockGetCrimeStatsUseCase struct {
	mock.Mock
}

func (m *MockGetCrimeStatsUseCase) Execute(ctx context.Context) (*entities.CrimeStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CrimeStats), args.Error(1)
}

// MockGetCrimeUseCase es un mock del caso de uso de obtención de un delito por ID
type MockGetCrimeUseCase struct {
	mock.Mock
}

func (m *MockGetCrimeUseCase) Execute(ctx context.Context, id string) (*entities.Crime, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Crime), args.Error(1)
}
