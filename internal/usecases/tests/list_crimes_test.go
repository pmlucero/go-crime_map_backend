package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-crime_map_backend/internal/domain/entities"
	domain_usecases "go-crime_map_backend/internal/domain/usecases"
	"go-crime_map_backend/internal/mocks"
	"go-crime_map_backend/internal/usecases"
	"go-crime_map_backend/internal/utils"
)

func TestListCrimesUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	startDate := now.Add(-24 * time.Hour)
	endDate := now

	t.Run("Listar todos los delitos", func(t *testing.T) {
		mockRepo := mocks.NewMockCrimeRepository()
		useCase := usecases.NewListCrimesUseCase(mockRepo)

		expectedCrimes := []entities.Crime{
			{
				ID:          "1",
				Title:       "Robo en tienda",
				Description: "Robo en tienda de conveniencia",
				Type:        string(entities.CrimeTypeRobo),
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		var nilTime *time.Time
		var nilString *string

		mockRepo.On("List", ctx, 1, 10, nilTime, nilTime, nilString, nilString).Return(expectedCrimes, int64(1), nil)

		output, err := useCase.Execute(ctx, domain_usecases.ListCrimesParams{
			Page:  1,
			Limit: 10,
		})

		assert.NoError(t, err)
		assert.Equal(t, expectedCrimes, output.Items)
		assert.Equal(t, int64(1), output.Total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error en rango de fechas inválido", func(t *testing.T) {
		mockRepo := mocks.NewMockCrimeRepository()
		useCase := usecases.NewListCrimesUseCase(mockRepo)

		output, err := useCase.Execute(ctx, domain_usecases.ListCrimesParams{
			Page:      1,
			Limit:     10,
			StartDate: &endDate,
			EndDate:   &startDate,
		})

		assert.Error(t, err)
		assert.Equal(t, usecases.ErrInvalidDateRange, err)
		assert.Nil(t, output)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Listar delitos con filtros", func(t *testing.T) {
		mockRepo := mocks.NewMockCrimeRepository()
		useCase := usecases.NewListCrimesUseCase(mockRepo)

		expectedCrimes := []entities.Crime{
			{
				ID:          "1",
				Title:       "Robo en tienda",
				Description: "Robo en tienda de conveniencia",
				Type:        string(entities.CrimeTypeRobo),
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		crimeType := string(entities.CrimeTypeRobo)
		status := string(entities.CrimeStatusActive)

		mockRepo.On("List", ctx, 1, 10, &startDate, &endDate, &crimeType, &status).Return(expectedCrimes, int64(1), nil)

		output, err := useCase.Execute(ctx, domain_usecases.ListCrimesParams{
			Page:      1,
			Limit:     10,
			Type:      &crimeType,
			Status:    &status,
			StartDate: &startDate,
			EndDate:   &endDate,
		})

		assert.NoError(t, err)
		assert.Equal(t, expectedCrimes, output.Items)
		assert.Equal(t, int64(1), output.Total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error al listar delitos", func(t *testing.T) {
		mockRepo := mocks.NewMockCrimeRepository()
		useCase := usecases.NewListCrimesUseCase(mockRepo)

		expectedError := errors.New("error al listar delitos")
		var nilTime *time.Time
		var nilString *string

		mockRepo.On("List", ctx, 1, 10, nilTime, nilTime, nilString, nilString).Return(nil, int64(0), expectedError)

		output, err := useCase.Execute(ctx, domain_usecases.ListCrimesParams{
			Page:  1,
			Limit: 10,
		})

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, output)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Página negativa", func(t *testing.T) {
		mockRepo := mocks.NewMockCrimeRepository()
		useCase := usecases.NewListCrimesUseCase(mockRepo)

		var nilTime *time.Time
		var nilString *string

		mockRepo.On("List", ctx, 1, 10, nilTime, nilTime, nilString, nilString).Return(nil, int64(0), nil)

		output, err := useCase.Execute(ctx, domain_usecases.ListCrimesParams{
			Page:  -1,
			Limit: 10,
		})

		assert.NoError(t, err)
		assert.NotNil(t, output)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Límite negativo", func(t *testing.T) {
		mockRepo := mocks.NewMockCrimeRepository()
		useCase := usecases.NewListCrimesUseCase(mockRepo)

		var nilTime *time.Time
		var nilString *string

		mockRepo.On("List", ctx, 1, 10, nilTime, nilTime, nilString, nilString).Return(nil, int64(0), nil)

		output, err := useCase.Execute(ctx, domain_usecases.ListCrimesParams{
			Page:  1,
			Limit: -5,
		})

		assert.NoError(t, err)
		assert.NotNil(t, output)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Límite muy grande", func(t *testing.T) {
		mockRepo := mocks.NewMockCrimeRepository()
		useCase := usecases.NewListCrimesUseCase(mockRepo)

		var nilTime *time.Time
		var nilString *string

		mockRepo.On("List", ctx, 1, 100, nilTime, nilTime, nilString, nilString).Return(nil, int64(0), nil)

		output, err := useCase.Execute(ctx, domain_usecases.ListCrimesParams{
			Page:  1,
			Limit: 1000,
		})

		assert.NoError(t, err)
		assert.NotNil(t, output)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Sin fechas", func(t *testing.T) {
		mockRepo := mocks.NewMockCrimeRepository()
		useCase := usecases.NewListCrimesUseCase(mockRepo)

		var nilTime *time.Time
		var nilString *string

		mockRepo.On("List", ctx, 1, 10, nilTime, nilTime, nilString, nilString).Return(nil, int64(0), nil)

		output, err := useCase.Execute(ctx, domain_usecases.ListCrimesParams{
			Page:  1,
			Limit: 10,
		})

		assert.NoError(t, err)
		assert.NotNil(t, output)
		mockRepo.AssertExpectations(t)
	})
}

func TestListCrimes_Execute_WithFilters(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := mocks.NewMockCrimeRepository()
	uc := usecases.NewListCrimesUseCase(repo)

	expectedCrimes := []entities.Crime{
		{
			ID:          "1",
			Title:       "Robo en tienda",
			Description: "Robo en tienda de conveniencia",
			Type:        string(entities.CrimeTypeRobo),
			Status:      string(entities.CrimeStatusActive),
		},
	}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	crimeType := string(entities.CrimeTypeRobo)
	status := string(entities.CrimeStatusActive)

	repo.On("List", ctx, 1, 10, &startDate, &endDate, &crimeType, &status).Return(expectedCrimes, int64(1), nil)

	// Act
	output, err := uc.Execute(ctx, domain_usecases.ListCrimesParams{
		Page:      1,
		Limit:     10,
		Type:      &crimeType,
		Status:    &status,
		StartDate: &startDate,
		EndDate:   &endDate,
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), output.Total)
	assert.Equal(t, expectedCrimes, output.Items)
	repo.AssertExpectations(t)
}

func TestListCrimes_Execute_WithTypeFilter(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := mocks.NewMockCrimeRepository()
	uc := usecases.NewListCrimesUseCase(repo)

	crimeType := string(entities.CrimeTypeRobo)
	var nilTime *time.Time
	var nilString *string

	repo.On("List", ctx, 1, 10, nilTime, nilTime, &crimeType, nilString).Return([]entities.Crime{}, int64(0), nil)

	// Act
	output, err := uc.Execute(ctx, domain_usecases.ListCrimesParams{
		Page:  1,
		Limit: 10,
		Type:  &crimeType,
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	repo.AssertExpectations(t)
}
