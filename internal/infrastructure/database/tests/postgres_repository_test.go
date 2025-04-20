package tests

import (
	"context"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
	"go-crime_map_backend/internal/infrastructure/database"
	"go-crime_map_backend/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *database.PostgresCrimeRepository {
	db := database.SetupTestDB(t)
	require.NotNil(t, db)

	repo := database.NewPostgresCrimeRepository(db.DB)
	require.NotNil(t, repo)

	t.Cleanup(func() {
		database.CleanupTestDB(t)
	})

	return repo
}

func TestPostgresCrimeRepository_Create(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	t.Run("Crear delito exitosamente", func(t *testing.T) {
		crime := &entities.Crime{
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
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, crime)
		assert.NoError(t, err)

		// Verificar que el delito fue creado
		savedCrime, err := repo.GetByID(ctx, crime.ID)
		assert.NoError(t, err)
		assert.NotNil(t, savedCrime)
		assert.Equal(t, crime.ID, savedCrime.ID)
		assert.Equal(t, crime.Title, savedCrime.Title)
		assert.Equal(t, crime.Description, savedCrime.Description)
		assert.Equal(t, crime.Type, savedCrime.Type)
		assert.Equal(t, crime.Status, savedCrime.Status)
		assert.Equal(t, crime.Location.Latitude, savedCrime.Location.Latitude)
		assert.Equal(t, crime.Location.Longitude, savedCrime.Location.Longitude)
		assert.Equal(t, crime.Location.Address, savedCrime.Location.Address)
	})

	t.Run("Error al crear delito - ID duplicado", func(t *testing.T) {
		crime := &entities.Crime{
			ID:          "test-id-2",
			Title:       "Test Crime",
			Description: "Test Description",
			Type:        "ROBBERY",
			Status:      "ACTIVE",
			Location: entities.Location{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Address:   "Test Address",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Primera creación
		err := repo.Create(ctx, crime)
		assert.NoError(t, err)

		// Intentar crear con el mismo ID
		err = repo.Create(ctx, crime)
		assert.Error(t, err)
	})
}

func TestPostgresCrimeRepository_GetByID(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	t.Run("Obtener delito existente", func(t *testing.T) {
		crime := &entities.Crime{
			ID:          "test-id-3",
			Title:       "Test Crime",
			Description: "Test Description",
			Type:        "ROBBERY",
			Status:      "ACTIVE",
			Location: entities.Location{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Address:   "Test Address",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, crime)
		assert.NoError(t, err)

		savedCrime, err := repo.GetByID(ctx, crime.ID)
		assert.NoError(t, err)
		assert.NotNil(t, savedCrime)
		assert.Equal(t, crime.ID, savedCrime.ID)
	})

	t.Run("Obtener delito inexistente", func(t *testing.T) {
		crime, err := repo.GetByID(ctx, "non-existent-id")
		assert.NoError(t, err)
		assert.Nil(t, crime)
	})
}

func TestPostgresCrimeRepository_Update(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	t.Run("Actualizar delito exitosamente", func(t *testing.T) {
		// Crear delito
		crime := &entities.Crime{
			ID:          "test-id-4",
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
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, crime)
		assert.NoError(t, err)

		// Actualizar delito
		crime.Title = "Robo en tienda actualizado"
		crime.Description = "Robo en tienda de conveniencia actualizado"
		crime.Status = string(entities.CrimeStatusInactive)
		crime.Location.Latitude = 40.7129
		crime.Location.Longitude = -74.0061
		crime.Location.Address = "Av. Santa Fe"
		crime.Location.AddressNumber = utils.StringPtr("5678")
		crime.UpdatedAt = time.Now()

		err = repo.Update(ctx, crime)
		assert.NoError(t, err)

		// Verificar actualización
		updatedCrime, err := repo.GetByID(ctx, crime.ID)
		assert.NoError(t, err)
		assert.NotNil(t, updatedCrime)
		assert.Equal(t, crime.Title, updatedCrime.Title)
		assert.Equal(t, crime.Description, updatedCrime.Description)
		assert.Equal(t, crime.Status, updatedCrime.Status)
		assert.Equal(t, crime.Location.Latitude, updatedCrime.Location.Latitude)
		assert.Equal(t, crime.Location.Longitude, updatedCrime.Location.Longitude)
		assert.Equal(t, crime.Location.Address, updatedCrime.Location.Address)
		assert.Equal(t, *crime.Location.AddressNumber, *updatedCrime.Location.AddressNumber)
	})

	t.Run("Error al actualizar delito inexistente", func(t *testing.T) {
		crime := &entities.Crime{
			ID:          "non-existent-id",
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
			UpdatedAt: time.Now(),
		}

		err := repo.Update(ctx, crime)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delito no encontrado")
	})

	t.Run("Error al actualizar delito eliminado", func(t *testing.T) {
		// Crear delito
		crime := &entities.Crime{
			ID:          "test-id-5",
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
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, crime)
		assert.NoError(t, err)

		// Eliminar delito
		err = repo.Delete(ctx, crime.ID)
		assert.NoError(t, err)

		// Intentar actualizar delito eliminado
		crime.Title = "Robo en tienda actualizado"
		crime.UpdatedAt = time.Now()
		err = repo.Update(ctx, crime)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delito no encontrado")
	})
}

func TestPostgresCrimeRepository_Delete(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	t.Run("Eliminar delito exitosamente", func(t *testing.T) {
		// Crear delito
		crime := &entities.Crime{
			ID:          "test-id-6",
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
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, crime)
		assert.NoError(t, err)

		// Eliminar delito
		err = repo.Delete(ctx, crime.ID)
		assert.NoError(t, err)

		// Verificar que el delito no se puede recuperar
		deletedCrime, err := repo.GetByID(ctx, crime.ID)
		assert.NoError(t, err)
		assert.Nil(t, deletedCrime)
	})

	t.Run("Error al eliminar delito inexistente", func(t *testing.T) {
		err := repo.Delete(ctx, "non-existent-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delito no encontrado")
	})

	t.Run("Error al eliminar delito ya eliminado", func(t *testing.T) {
		// Crear delito
		crime := &entities.Crime{
			ID:          "test-id-7",
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
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, crime)
		assert.NoError(t, err)

		// Eliminar delito
		err = repo.Delete(ctx, crime.ID)
		assert.NoError(t, err)

		// Intentar eliminar el mismo delito nuevamente
		err = repo.Delete(ctx, crime.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delito no encontrado")
	})
}

func TestPostgresCrimeRepository_GetAll(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	t.Run("Obtener todos los delitos exitosamente", func(t *testing.T) {
		// Crear algunos delitos de prueba
		crime1 := &entities.Crime{
			ID:          "test-id-getall-1",
			Title:       "Robo en tienda 1",
			Description: "Robo en tienda de conveniencia 1",
			Type:        string(entities.CrimeTypeRobo),
			Status:      string(entities.CrimeStatusActive),
			Location: entities.Location{
				Latitude:      40.7128,
				Longitude:     -74.0060,
				Address:       "Av. Corrientes",
				AddressNumber: utils.StringPtr("1234"),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		crime2 := &entities.Crime{
			ID:          "test-id-getall-2",
			Title:       "Robo en tienda 2",
			Description: "Robo en tienda de conveniencia 2",
			Type:        string(entities.CrimeTypeRobo),
			Status:      string(entities.CrimeStatusActive),
			Location: entities.Location{
				Latitude:      40.7129,
				Longitude:     -74.0061,
				Address:       "Av. Santa Fe",
				AddressNumber: utils.StringPtr("5678"),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, crime1)
		assert.NoError(t, err)
		err = repo.Create(ctx, crime2)
		assert.NoError(t, err)

		// Obtener todos los delitos
		crimes, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, crimes, 2)
	})

	t.Run("Error al ejecutar la consulta", func(t *testing.T) {
		// Cerrar la conexión para forzar un error
		repo.DB.Close()
		_, err := repo.GetAll(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error al obtener delitos")
	})
}

func TestPostgresCrimeRepository_List(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	// Crear algunos delitos de prueba
	crimes := []*entities.Crime{
		{
			ID:          "test-id-6",
			Title:       "Robo 1",
			Description: "Descripción 1",
			Type:        "ROBBERY",
			Status:      "ACTIVE",
			Location: entities.Location{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Address:   "Address 1",
			},
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(),
		},
		{
			ID:          "test-id-7",
			Title:       "Robo 2",
			Description: "Descripción 2",
			Type:        "ROBBERY",
			Status:      "INACTIVE",
			Location: entities.Location{
				Latitude:  40.7129,
				Longitude: -74.0061,
				Address:   "Address 2",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:          "test-id-8",
			Title:       "Asalto 1",
			Description: "Descripción 3",
			Type:        "ASSAULT",
			Status:      "ACTIVE",
			Location: entities.Location{
				Latitude:  40.7130,
				Longitude: -74.0062,
				Address:   "Address 3",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, crime := range crimes {
		err := repo.Create(ctx, crime)
		assert.NoError(t, err)
	}

	t.Run("Listar todos los delitos", func(t *testing.T) {
		filter := repositories.ListCrimesFilter{
			Page:  1,
			Limit: 10,
		}

		result, err := repo.List(ctx, filter)
		assert.NoError(t, err)
		assert.Len(t, result, 3)
	})

	t.Run("Filtrar por tipo", func(t *testing.T) {
		filter := repositories.ListCrimesFilter{
			Type:  "ROBBERY",
			Page:  1,
			Limit: 10,
		}

		result, err := repo.List(ctx, filter)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		for _, crime := range result {
			assert.Equal(t, "ROBBERY", crime.Type)
		}
	})

	t.Run("Filtrar por estado", func(t *testing.T) {
		filter := repositories.ListCrimesFilter{
			Status: "ACTIVE",
			Page:   1,
			Limit:  10,
		}

		result, err := repo.List(ctx, filter)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		for _, crime := range result {
			assert.Equal(t, "ACTIVE", crime.Status)
		}
	})

	t.Run("Filtrar por fecha", func(t *testing.T) {
		filter := repositories.ListCrimesFilter{
			StartDate: time.Now().Add(-12 * time.Hour),
			Page:      1,
			Limit:     10,
		}

		result, err := repo.List(ctx, filter)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("Filtrar por ubicación", func(t *testing.T) {
		filter := repositories.ListCrimesFilter{
			Latitude:  40.7128,
			Longitude: -74.0060,
			RadiusKm:  1,
			Page:      1,
			Limit:     10,
		}

		result, err := repo.List(ctx, filter)
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("Paginación", func(t *testing.T) {
		filter := repositories.ListCrimesFilter{
			Page:  1,
			Limit: 2,
		}

		result, err := repo.List(ctx, filter)
		assert.NoError(t, err)
		assert.Len(t, result, 2)

		filter.Page = 2
		result, err = repo.List(ctx, filter)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
	})
}

func TestPostgresCrimeRepository_Update_Errors(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	t.Run("Error al ejecutar la consulta", func(t *testing.T) {
		crime := &entities.Crime{
			ID:          "test-id-update-error",
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
			UpdatedAt: time.Now(),
		}

		// Cerrar la conexión para forzar un error
		repo.DB.Close()
		err := repo.Update(ctx, crime)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error al actualizar delito")
	})
}

func TestPostgresCrimeRepository_Delete_Errors(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	t.Run("Error al ejecutar la consulta", func(t *testing.T) {
		// Cerrar la conexión para forzar un error
		repo.DB.Close()
		err := repo.Delete(ctx, "test-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error al eliminar delito")
	})

	t.Run("Delito no encontrado", func(t *testing.T) {
		repo := setupTestDB(t)
		err := repo.Delete(ctx, "non-existent-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delito no encontrado")
	})
}
