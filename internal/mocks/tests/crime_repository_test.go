package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/mocks"
	"go-crime_map_backend/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestMockCrimeRepository_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("Creación exitosa", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		crime := &entities.Crime{
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
		}

		mockRepo.On("Create", ctx, crime).Return(nil).Once()

		err := mockRepo.Create(ctx, crime)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error al crear", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		crime := &entities.Crime{
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
		}
		expectedError := errors.New("error al crear el delito")

		mockRepo.On("Create", ctx, crime).Return(expectedError).Once()

		err := mockRepo.Create(ctx, crime)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestMockCrimeRepository_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("Obtención exitosa", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		expectedCrime := &entities.Crime{
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
		}

		mockRepo.On("GetByID", ctx, "1").Return(expectedCrime, nil).Once()

		result, err := mockRepo.GetByID(ctx, "1")

		assert.NoError(t, err)
		assert.Equal(t, expectedCrime, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error al obtener", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		expectedError := errors.New("error al obtener el delito")

		mockRepo.On("GetByID", ctx, "1").Return(nil, expectedError).Once()

		result, err := mockRepo.GetByID(ctx, "1")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestMockCrimeRepository_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("Obtención exitosa", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		expectedCrimes := []*entities.Crime{
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
			},
			{
				ID:          "2",
				Title:       "Asalto",
				Description: "Asalto a mano armada",
				Type:        "ASSAULT",
				Status:      "ACTIVE",
			},
		}

		mockRepo.On("GetAll", ctx).Return(expectedCrimes, nil).Once()

		result, err := mockRepo.GetAll(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedCrimes, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error al obtener todos", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		expectedError := errors.New("error al obtener los delitos")

		mockRepo.On("GetAll", ctx).Return(nil, expectedError).Once()

		result, err := mockRepo.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestMockCrimeRepository_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("Actualización exitosa", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		crime := &entities.Crime{
			ID:          "1",
			Title:       "Robo en tienda actualizado",
			Description: "Robo en tienda de conveniencia actualizado",
			Type:        string(entities.CrimeTypeRobo),
			Status:      string(entities.CrimeStatusActive),
			Location: entities.Location{
				Latitude:      40.7128,
				Longitude:     -74.0060,
				Address:       "Av. Corrientes",
				AddressNumber: utils.StringPtr("1234"),
			},
		}

		mockRepo.On("Update", ctx, crime).Return(nil).Once()

		err := mockRepo.Update(ctx, crime)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error al actualizar", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		crime := &entities.Crime{
			ID:          "1",
			Title:       "Robo en tienda actualizado",
			Description: "Robo en tienda de conveniencia actualizado",
			Type:        string(entities.CrimeTypeRobo),
			Status:      string(entities.CrimeStatusActive),
			Location: entities.Location{
				Latitude:      40.7128,
				Longitude:     -74.0060,
				Address:       "Av. Corrientes",
				AddressNumber: utils.StringPtr("1234"),
			},
		}
		expectedError := errors.New("error al actualizar el delito")

		mockRepo.On("Update", ctx, crime).Return(expectedError).Once()

		err := mockRepo.Update(ctx, crime)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestMockCrimeRepository_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("Eliminación exitosa", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		id := "1"

		mockRepo.On("Delete", ctx, id).Return(nil).Once()

		err := mockRepo.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error al eliminar", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		id := "1"
		expectedError := errors.New("error al eliminar el delito")

		mockRepo.On("Delete", ctx, id).Return(expectedError).Once()

		err := mockRepo.Delete(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestMockCrimeRepository_List(t *testing.T) {
	ctx := context.Background()

	t.Run("Listado exitoso", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
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
			},
			{
				ID:          "2",
				Title:       "Agresión en calle",
				Description: "Agresión en la vía pública",
				Type:        string(entities.CrimeTypeAgresion),
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
			},
		}
		expectedTotal := int64(2)
		page := 1
		limit := 10
		startDate := time.Now().Add(-24 * time.Hour)
		endDate := time.Now()
		crimeType := string(entities.CrimeTypeRobo)
		status := string(entities.CrimeStatusActive)

		mockRepo.On("List", ctx, page, limit, &startDate, &endDate, &crimeType, &status).
			Return(expectedCrimes, expectedTotal, nil).Once()

		result, total, err := mockRepo.List(ctx, page, limit, &startDate, &endDate, &crimeType, &status)

		assert.NoError(t, err)
		assert.Equal(t, expectedCrimes, result)
		assert.Equal(t, expectedTotal, total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error al listar", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		expectedError := errors.New("error al listar los delitos")
		page := 1
		limit := 10
		var startDate, endDate *time.Time
		var crimeType, status *string

		mockRepo.On("List", ctx, page, limit, startDate, endDate, crimeType, status).
			Return(nil, int64(0), expectedError).Once()

		result, total, err := mockRepo.List(ctx, page, limit, startDate, endDate, crimeType, status)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestMockCrimeRepository_GetStats(t *testing.T) {
	ctx := context.Background()

	t.Run("Obtención exitosa", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		expectedStats := &entities.CrimeStats{
			TotalCrimes:    100,
			ActiveCrimes:   70,
			InactiveCrimes: 30,
			CrimesByType: map[string]int64{
				string(entities.CrimeTypeRobo):     50,
				string(entities.CrimeTypeAgresion): 30,
				string(entities.CrimeTypeHurto):    20,
			},
			CrimesByStatus: map[string]int64{
				string(entities.CrimeStatusActive):   70,
				string(entities.CrimeStatusInactive): 30,
			},
			LastUpdate: time.Date(2025, 4, 15, 21, 42, 14, 862708000, time.Local),
		}

		mockRepo.On("GetStats", ctx).Return(expectedStats, nil).Once()

		result, err := mockRepo.GetStats(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedStats, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error al obtener estadísticas", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		expectedError := errors.New("error al obtener estadísticas")

		mockRepo.On("GetStats", ctx).Return(nil, expectedError).Once()

		result, err := mockRepo.GetStats(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}
