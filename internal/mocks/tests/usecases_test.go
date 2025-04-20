package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/usecases"
	"go-crime_map_backend/internal/mocks"

	"github.com/stretchr/testify/assert"
)

func TestMockCreateCrimeUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	// Caso exitoso
	t.Run("Creación exitosa", func(t *testing.T) {
		mockUseCase := new(mocks.MockCreateCrimeUseCase)
		expectedCrime := &entities.Crime{
			ID:          "1",
			Title:       "Robo",
			Description: "Robo en tienda",
			Type:        "ROBBERY",
			Status:      "ACTIVE",
			Location: entities.Location{
				Latitude:  -34.603722,
				Longitude: -58.381592,
			},
		}

		input := usecases.CreateCrimeInput{
			Title:       "Robo",
			Description: "Robo en tienda",
			Type:        "ROBBERY",
			Latitude:    -34.603722,
			Longitude:   -58.381592,
		}

		mockUseCase.On("Execute", ctx, input).Return(expectedCrime, nil).Once()

		result, err := mockUseCase.Execute(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, expectedCrime, result)
		mockUseCase.AssertExpectations(t)
	})

	// Caso de error
	t.Run("Error al crear", func(t *testing.T) {
		mockUseCase := new(mocks.MockCreateCrimeUseCase)
		expectedError := errors.New("error al crear el delito")
		input := usecases.CreateCrimeInput{
			Title:       "Robo",
			Description: "Robo en tienda",
			Type:        "ROBBERY",
			Latitude:    -34.603722,
			Longitude:   -58.381592,
		}

		mockUseCase.On("Execute", ctx, input).Return(nil, expectedError).Once()

		result, err := mockUseCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockUseCase.AssertExpectations(t)
	})
}

func TestMockListCrimesUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	// Caso exitoso
	t.Run("Listado exitoso", func(t *testing.T) {
		mockUseCase := new(mocks.MockListCrimesUseCase)
		expectedCrimes := &entities.CrimeList{
			Total: 2,
			Items: []entities.Crime{
				{
					ID:          "1",
					Title:       "Robo",
					Description: "Robo en tienda",
					Type:        "ROBBERY",
					Status:      "ACTIVE",
				},
				{
					ID:          "2",
					Title:       "Asalto",
					Description: "Asalto a mano armada",
					Type:        "ASSAULT",
					Status:      "ACTIVE",
				},
			},
		}

		params := usecases.ListCrimesParams{
			Page:  1,
			Limit: 10,
		}

		mockUseCase.On("Execute", ctx, params).Return(expectedCrimes, nil).Once()

		result, err := mockUseCase.Execute(ctx, params)

		assert.NoError(t, err)
		assert.Equal(t, expectedCrimes, result)
		mockUseCase.AssertExpectations(t)
	})

	// Caso de error
	t.Run("Error al listar", func(t *testing.T) {
		mockUseCase := new(mocks.MockListCrimesUseCase)
		expectedError := errors.New("error al listar delitos")
		params := usecases.ListCrimesParams{
			Page:  1,
			Limit: 10,
		}

		mockUseCase.On("Execute", ctx, params).Return(nil, expectedError).Once()

		result, err := mockUseCase.Execute(ctx, params)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockUseCase.AssertExpectations(t)
	})
}

func TestMockUpdateCrimeStatusUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	// Caso exitoso
	t.Run("Actualización exitosa", func(t *testing.T) {
		mockUseCase := new(mocks.MockUpdateCrimeStatusUseCase)
		input := usecases.UpdateCrimeStatusInput{
			ID:     "1",
			Status: "INACTIVE",
		}

		mockUseCase.On("Execute", ctx, input).Return(nil).Once()

		err := mockUseCase.Execute(ctx, input)

		assert.NoError(t, err)
		mockUseCase.AssertExpectations(t)
	})

	// Caso de error
	t.Run("Error al actualizar", func(t *testing.T) {
		mockUseCase := new(mocks.MockUpdateCrimeStatusUseCase)
		expectedError := errors.New("error al actualizar estado")
		input := usecases.UpdateCrimeStatusInput{
			ID:     "1",
			Status: "INACTIVE",
		}

		mockUseCase.On("Execute", ctx, input).Return(expectedError).Once()

		err := mockUseCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockUseCase.AssertExpectations(t)
	})
}

func TestMockDeleteCrimeUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	// Caso exitoso
	t.Run("Eliminación exitosa", func(t *testing.T) {
		mockUseCase := new(mocks.MockDeleteCrimeUseCase)
		id := "1"

		mockUseCase.On("Execute", ctx, id).Return(nil).Once()

		err := mockUseCase.Execute(ctx, id)

		assert.NoError(t, err)
		mockUseCase.AssertExpectations(t)
	})

	// Caso de error
	t.Run("Error al eliminar", func(t *testing.T) {
		mockUseCase := new(mocks.MockDeleteCrimeUseCase)
		expectedError := errors.New("error al eliminar delito")
		id := "1"

		mockUseCase.On("Execute", ctx, id).Return(expectedError).Once()

		err := mockUseCase.Execute(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockUseCase.AssertExpectations(t)
	})
}

func TestMockGetCrimeStatsUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	// Caso exitoso
	t.Run("Obtención exitosa", func(t *testing.T) {
		mockUseCase := new(mocks.MockGetCrimeStatsUseCase)
		expectedStats := &entities.CrimeStats{
			TotalCrimes:    100,
			ActiveCrimes:   80,
			InactiveCrimes: 20,
			CrimesByType: map[string]int64{
				"ROBBERY": 50,
				"ASSAULT": 30,
			},
			LastUpdate: time.Date(2025, 4, 15, 21, 42, 14, 862708000, time.Local),
		}

		mockUseCase.On("Execute", ctx).Return(expectedStats, nil).Once()

		result, err := mockUseCase.Execute(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedStats, result)
		mockUseCase.AssertExpectations(t)
	})

	// Caso de error
	t.Run("Error al obtener estadísticas", func(t *testing.T) {
		mockUseCase := new(mocks.MockGetCrimeStatsUseCase)
		expectedError := errors.New("error al obtener estadísticas")

		mockUseCase.On("Execute", ctx).Return(nil, expectedError).Once()

		result, err := mockUseCase.Execute(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockUseCase.AssertExpectations(t)
	})
}

func TestMockGetCrimeUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	// Caso exitoso
	t.Run("Obtención exitosa", func(t *testing.T) {
		mockUseCase := new(mocks.MockGetCrimeUseCase)
		expectedCrime := &entities.Crime{
			ID:          "1",
			Title:       "Robo",
			Description: "Robo en tienda",
			Type:        "ROBBERY",
			Status:      "ACTIVE",
		}

		id := "1"

		mockUseCase.On("Execute", ctx, id).Return(expectedCrime, nil).Once()

		result, err := mockUseCase.Execute(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedCrime, result)
		mockUseCase.AssertExpectations(t)
	})

	// Caso de error
	t.Run("Error al obtener delito", func(t *testing.T) {
		mockUseCase := new(mocks.MockGetCrimeUseCase)
		expectedError := errors.New("error al obtener delito")
		id := "1"

		mockUseCase.On("Execute", ctx, id).Return(nil, expectedError).Once()

		result, err := mockUseCase.Execute(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockUseCase.AssertExpectations(t)
	})
}
