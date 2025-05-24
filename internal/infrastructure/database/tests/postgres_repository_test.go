package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
	"go-crime_map_backend/internal/domain/usecases"
	"go-crime_map_backend/internal/domain/usecases/mocks"
	"go-crime_map_backend/internal/infrastructure/database"
	"go-crime_map_backend/internal/interface/controllers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// GenerateUUID genera un UUID único
func GenerateUUID() string {
	return uuid.New().String()
}

func SetupTestDB(t *testing.T) {
	db := database.SetupTestDB(t)
	require.NotNil(t, db)

	// Usar db.DB para ejecutar las queries
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS crimes (
			id BIGSERIAL PRIMARY KEY,
			uuid VARCHAR(36) NOT NULL UNIQUE,
			title VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			crime_type VARCHAR(50) NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			status VARCHAR(20) NOT NULL,
			address VARCHAR(255) NOT NULL,
			address_number VARCHAR(20) NOT NULL,
			city VARCHAR(100) NOT NULL,
			province VARCHAR(100) NOT NULL,
			country VARCHAR(100) NOT NULL,
			zip_code VARCHAR(20) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
			deleted_at TIMESTAMP WITH TIME ZONE
		)
	`)
	assert.NoError(t, err)

	// Luego crear los índices
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_crimes_location 
		ON crimes USING GIST (ll_to_earth(latitude, longitude))
	`)
	assert.NoError(t, err)
}

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

// Helper function para crear una Location válida
func createValidLocation() entities.Location {
	return entities.Location{
		Latitude:      40.7128,
		Longitude:     -74.0060,
		Address:       "Av. Corrientes",
		AddressNumber: "1234",
		City:          "Buenos Aires",
		Province:      "CABA",
		Country:       "Argentina",
		ZipCode:       "C1043AAZ",
	}
}

func TestCreate_ValidCrime_Success(t *testing.T) {
	// Arrange
	repo := setupTestDB(t)
	ctx := context.Background()
	uuid := GenerateUUID()
	crime := &entities.Crime{
		UUID:        uuid,
		Title:       "Robo en tienda",
		Description: "Robo en tienda de conveniencia",
		Type:        string(entities.CrimeTypeRobo),
		Status:      string(entities.CrimeStatusActive),
		Location:    createValidLocation(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Act
	err := repo.Create(ctx, crime)

	// Assert
	assert.NoError(t, err)
	savedCrime, err := repo.GetByUUID(ctx, uuid)
	assert.NoError(t, err)
	assert.NotNil(t, savedCrime)
	assert.Equal(t, crime.UUID, savedCrime.UUID)
	assert.Equal(t, crime.Title, savedCrime.Title)
	assert.Equal(t, crime.Description, savedCrime.Description)
	assert.Equal(t, crime.Type, savedCrime.Type)
	assert.Equal(t, crime.Status, savedCrime.Status)
	assert.Equal(t, crime.Location.Latitude, savedCrime.Location.Latitude)
	assert.Equal(t, crime.Location.Longitude, savedCrime.Location.Longitude)
	assert.Equal(t, crime.Location.Address, savedCrime.Location.Address)
}

func TestCreate_DuplicateID_ReturnsError(t *testing.T) {
	// Arrange
	repo := setupTestDB(t)
	ctx := context.Background()
	uuid := GenerateUUID()
	crime := &entities.Crime{
		UUID:        uuid,
		Title:       "Test Crime",
		Description: "Test Description",
		Type:        string(entities.CrimeTypeRobo),
		Status:      string(entities.CrimeStatusActive),
		Location:    createValidLocation(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Act
	err := repo.Create(ctx, crime)
	assert.NoError(t, err)

	// Assert
	duplicateCrime := &entities.Crime{
		UUID:        uuid,
		Title:       "Test Crime",
		Description: "Test Description",
		Type:        string(entities.CrimeTypeRobo),
		Status:      string(entities.CrimeStatusActive),
		Location:    createValidLocation(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = repo.Create(ctx, duplicateCrime)
	assert.Error(t, err)
}

func TestGetByUUID_ExistingCrime_ReturnsSuccess(t *testing.T) {
	// Arrange
	repo := setupTestDB(t)
	ctx := context.Background()
	uuid := GenerateUUID()
	crime := &entities.Crime{
		UUID:        uuid,
		Title:       "Test Crime",
		Description: "Test Description",
		Type:        string(entities.CrimeTypeRobo),
		Status:      string(entities.CrimeStatusActive),
		Location:    createValidLocation(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, crime)
	assert.NoError(t, err)

	// Act
	savedCrime, err := repo.GetByUUID(ctx, uuid)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, savedCrime)
	assert.Equal(t, crime.UUID, savedCrime.UUID)
}

func TestGetByUUID_NonExistentCrime_ReturnsNil(t *testing.T) {
	// Arrange
	repo := setupTestDB(t)
	ctx := context.Background()

	// Act
	crime, err := repo.GetByUUID(ctx, GenerateUUID())

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, crime)
}

func TestUpdate_ValidCrime_Success(t *testing.T) {
	// Arrange
	repo := setupTestDB(t)
	ctx := context.Background()
	uuid := GenerateUUID()
	crime := &entities.Crime{
		UUID:        uuid,
		Title:       "Robo en tienda",
		Description: "Robo en tienda de conveniencia",
		Type:        string(entities.CrimeTypeRobo),
		Status:      string(entities.CrimeStatusActive),
		Location:    createValidLocation(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, crime)
	assert.NoError(t, err)

	// Act
	crime.Title = "Robo en tienda actualizado"
	crime.Description = "Robo en tienda de conveniencia actualizado"
	crime.Status = string(entities.CrimeStatusInactive)
	newLocation := createValidLocation()
	newLocation.Address = "Av. Santa Fe"
	newLocation.AddressNumber = "5678"
	crime.Location = newLocation
	crime.UpdatedAt = time.Now()

	err = repo.Update(ctx, crime)

	// Assert
	assert.NoError(t, err)
	updatedCrime, err := repo.GetByUUID(ctx, uuid)
	assert.NoError(t, err)
	assert.NotNil(t, updatedCrime)
	assert.Equal(t, crime.Title, updatedCrime.Title)
	assert.Equal(t, crime.Description, updatedCrime.Description)
	assert.Equal(t, crime.Status, updatedCrime.Status)
	assert.Equal(t, crime.Location.Address, updatedCrime.Location.Address)
	assert.Equal(t, crime.Location.AddressNumber, updatedCrime.Location.AddressNumber)
}

func TestUpdate_NonExistentCrime_ReturnsError(t *testing.T) {
	// Arrange
	repo := setupTestDB(t)
	ctx := context.Background()
	crime := &entities.Crime{
		UUID:        GenerateUUID(),
		Title:       "Robo en tienda",
		Description: "Robo en tienda de conveniencia",
		Type:        string(entities.CrimeTypeRobo),
		Status:      string(entities.CrimeStatusActive),
		Location:    createValidLocation(),
		UpdatedAt:   time.Now(),
	}

	// Act
	err := repo.Update(ctx, crime)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delito no encontrado")
}

func TestDelete_ExistingCrime_Success(t *testing.T) {
	// Arrange
	repo := setupTestDB(t)
	ctx := context.Background()
	uuid := GenerateUUID()
	crime := &entities.Crime{
		UUID:        uuid,
		Description: "Test Description",
		Type:        string(entities.CrimeTypeRobo),
		Status:      string(entities.CrimeStatusActive),
		Location:    createValidLocation(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, crime)
	assert.NoError(t, err)
	assert.NotZero(t, crime.ID)

	// Act
	err = repo.Delete(ctx, crime.ID)

	// Assert
	assert.NoError(t, err)
	deletedCrime, err := repo.GetByUUID(ctx, uuid)
	assert.Nil(t, deletedCrime)
	assert.NoError(t, err)
}

func TestDelete_NonExistentCrime_ReturnsError(t *testing.T) {
	// Arrange
	repo := setupTestDB(t)
	ctx := context.Background()

	// Act
	err := repo.Delete(ctx, 999999)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delito no encontrado")
}

func TestDelete_AlreadyDeletedCrime_ReturnsError(t *testing.T) {
	// Arrange
	repo := setupTestDB(t)
	ctx := context.Background()
	uuid := GenerateUUID()
	crime := &entities.Crime{
		UUID:        uuid,
		Title:       "Test Crime",
		Description: "Test Description",
		Type:        "ROBBERY",
		Status:      "ACTIVE",
		Location:    createValidLocation(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, crime)
	assert.NoError(t, err)

	// Primera eliminación
	err = repo.Delete(ctx, crime.ID)
	assert.NoError(t, err)

	// Segunda eliminación
	err = repo.Delete(ctx, crime.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delito no encontrado")
}

func TestGetAll_WithCrimes_ReturnsAllCrimes(t *testing.T) {
	// Arrange
	repo := setupTestDB(t)
	ctx := context.Background()
	crime1 := &entities.Crime{
		UUID:        GenerateUUID(),
		Title:       "Test Crime 1",
		Description: "Test Description 1",
		Type:        "ROBBERY",
		Status:      "ACTIVE",
		Location:    createValidLocation(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	crime2 := &entities.Crime{
		UUID:        GenerateUUID(),
		Title:       "Test Crime 2",
		Description: "Test Description 2",
		Type:        "ROBBERY",
		Status:      "ACTIVE",
		Location:    createValidLocation(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, crime1)
	assert.NoError(t, err)
	err = repo.Create(ctx, crime2)
	assert.NoError(t, err)

	// Act
	crimes, err := repo.GetAll(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, crimes, 2)
}

func TestList_WithFilters_ReturnsFilteredCrimes(t *testing.T) {
	// Arrange
	repo := setupTestDB(t)
	ctx := context.Background()
	now := time.Now()
	crime1 := &entities.Crime{
		UUID:        GenerateUUID(),
		Title:       "Test Crime 1",
		Description: "Test Description 1",
		Type:        "ROBBERY",
		Status:      "ACTIVE",
		Location:    createValidLocation(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	crime2 := &entities.Crime{
		UUID:        GenerateUUID(),
		Title:       "Test Crime 2",
		Description: "Test Description 2",
		Type:        "ASSAULT",
		Status:      "INACTIVE",
		Location:    createValidLocation(),
		CreatedAt:   now.Add(time.Hour),
		UpdatedAt:   now.Add(time.Hour),
	}
	err := repo.Create(ctx, crime1)
	assert.NoError(t, err)
	err = repo.Create(ctx, crime2)
	assert.NoError(t, err)

	// Act
	filters := repositories.ListCrimesFilter{
		Type:      "ROBBERY",
		Status:    "ACTIVE",
		StartDate: now.Add(-time.Hour),
		EndDate:   now.Add(time.Hour),
		Latitude:  40.7128,
		Longitude: -74.0060,
		RadiusKm:  1,
		Page:      1,
		Limit:     10,
	}
	crimes, err := repo.List(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, crimes, 1)
	assert.Equal(t, crime1.UUID, crimes[0].UUID)
}

func TestList_WithPagination_ReturnsPagedCrimes(t *testing.T) {
	// Arrange
	repo := setupTestDB(t)
	ctx := context.Background()
	now := time.Now()
	for i := 0; i < 5; i++ {
		crime := &entities.Crime{
			UUID:        GenerateUUID(),
			Title:       "Test Crime",
			Description: "Test Description",
			Type:        "ROBBERY",
			Status:      "ACTIVE",
			Location:    createValidLocation(),
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		err := repo.Create(ctx, crime)
		assert.NoError(t, err)
	}

	// Act
	filters := repositories.ListCrimesFilter{
		Page:  1,
		Limit: 2,
	}
	crimes, err := repo.List(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, crimes, 2)
}

func TestUpdate_DatabaseError_ReturnsError(t *testing.T) {
	// Arrange
	repo := setupTestDB(t)
	ctx := context.Background()
	crime := &entities.Crime{
		UUID:        GenerateUUID(),
		Title:       "Test Crime",
		Description: "Test Description",
		Type:        "ROBBERY",
		Status:      "ACTIVE",
		Location:    createValidLocation(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Act
	err := repo.Update(ctx, crime)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delito no encontrado")
}

func TestDelete_DatabaseError_ReturnsError(t *testing.T) {
	// Arrange
	repo := setupTestDB(t)
	ctx := context.Background()

	// Act
	err := repo.Delete(ctx, 999999)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delito no encontrado")
}

func TestUpdateCrimeStatus_Success(t *testing.T) {
	mockUseCase := &mocks.UpdateCrimeStatusUseCase{}

	mockUseCase.On("Execute", mock.Anything, usecases.UpdateCrimeStatusInput{
		UUID:   "test-uuid",
		Status: "INACTIVE",
	}).Return(nil)

	controller := controllers.NewCrimeController(
		nil,
		nil,
		mockUseCase,
		nil,
		nil,
		nil,
	)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "test-uuid"}}
	c.Request = httptest.NewRequest(
		"PATCH",
		"/crimes/test-uuid/status",
		strings.NewReader(`{"status":"INACTIVE"}`),
	)
	c.Request.Header.Set("Content-Type", "application/json")

	controller.UpdateCrimeStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}
