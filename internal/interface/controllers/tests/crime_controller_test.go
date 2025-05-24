package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/usecases"
	"go-crime_map_backend/internal/domain/usecases/mocks"
	"go-crime_map_backend/internal/interface/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Definición de errores
var (
	ErrInvalidCrimeStatus = errors.New("invalid crime status")
	ErrCrimeNotFound      = errors.New("crime not found")
	ErrInvalidAPIKey      = errors.New("invalid API key")
)

// Definición de tipos de entrada para los casos de uso
type DeleteCrimeInput struct {
	ID string
}

type GetCrimeInput struct {
	ID string
}

// MockCreateCrimeUseCase es un mock para el caso de uso de creación de delitos
type MockCreateCrimeUseCase struct {
	mock.Mock
}

func (m *MockCreateCrimeUseCase) Execute(ctx context.Context, input usecases.CreateCrimeInput) (*entities.Crime, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*entities.Crime), args.Error(1)
}

// MockListCrimesUseCase es un mock para el caso de uso de listado de delitos
type MockListCrimesUseCase struct {
	mock.Mock
}

func (m *MockListCrimesUseCase) Execute(ctx context.Context, params usecases.ListCrimesParams) (*entities.CrimeList, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*entities.CrimeList), args.Error(1)
}

// MockUpdateCrimeStatusUseCase es un mock para el caso de uso de actualización de estado
type MockUpdateCrimeStatusUseCase struct {
	mock.Mock
}

func (m *MockUpdateCrimeStatusUseCase) Execute(ctx context.Context, input usecases.UpdateCrimeStatusInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

// MockDeleteCrimeUseCase es un mock para el caso de uso de eliminación de delitos
type MockDeleteCrimeUseCase struct {
	mock.Mock
}

func (m *MockDeleteCrimeUseCase) Execute(ctx context.Context, crimeID string) error {
	args := m.Called(ctx, crimeID)
	return args.Error(0)
}

// MockGetCrimeStatsUseCase es un mock para el caso de uso de obtención de estadísticas
type MockGetCrimeStatsUseCase struct {
	mock.Mock
}

func (m *MockGetCrimeStatsUseCase) Execute(ctx context.Context) (*entities.CrimeStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(*entities.CrimeStats), args.Error(1)
}

// MockGetCrimeUseCase es un mock para el caso de uso de obtención de un delito
type MockGetCrimeUseCase struct {
	mock.Mock
}

func (m *MockGetCrimeUseCase) Execute(ctx context.Context, crimeID string) (*entities.Crime, error) {
	args := m.Called(ctx, crimeID)
	return args.Get(0).(*entities.Crime), args.Error(1)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestCrimeController_CreateCrime(t *testing.T) {
	// Arrange
	mockCreateUseCase := new(mocks.CreateCrimeUseCase)
	controller := controllers.NewCrimeController(
		mockCreateUseCase,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	router := setupRouter()
	router.POST("/crimes", controller.CreateCrime)

	addressNumber := "123"
	city := "Buenos Aires"
	province := "Buenos Aires"
	country := "Argentina"
	zipCode := "1000"

	reqBody := controllers.CreateCrimeRequest{
		Title:         "Test Crime",
		Description:   "Test Description",
		Type:          "ROBO",
		Latitude:      -34.603722,
		Longitude:     -58.381592,
		Address:       "Test Address",
		AddressNumber: addressNumber,
		City:          city,
		Province:      province,
		Country:       country,
		ZipCode:       zipCode,
	}

	expectedCrime := &entities.Crime{
		ID:          1,
		Title:       reqBody.Title,
		Description: reqBody.Description,
		Type:        "ROBO",
		Status:      "ACTIVE",
		Location: entities.Location{
			Latitude:      reqBody.Latitude,
			Longitude:     reqBody.Longitude,
			Address:       reqBody.Address,
			AddressNumber: reqBody.AddressNumber,
			City:          reqBody.City,
			Province:      reqBody.Province,
			Country:       reqBody.Country,
			ZipCode:       reqBody.ZipCode,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	input := usecases.CreateCrimeInput{
		Title:         reqBody.Title,
		Description:   reqBody.Description,
		Type:          reqBody.Type,
		Latitude:      reqBody.Latitude,
		Longitude:     reqBody.Longitude,
		Address:       reqBody.Address,
		AddressNumber: reqBody.AddressNumber,
		City:          reqBody.City,
		Province:      reqBody.Province,
		Country:       reqBody.Country,
		ZipCode:       reqBody.ZipCode,
	}

	mockCreateUseCase.On("Execute", mock.Anything, input).Return(expectedCrime, nil)

	// Act
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/crimes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "valid-key")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	mockCreateUseCase.AssertExpectations(t)
}

func TestCrimeController_ListCrimes(t *testing.T) {
	// Arrange
	mockListUseCase := new(mocks.ListCrimesUseCase)
	controller := controllers.NewCrimeController(
		nil,
		mockListUseCase,
		nil,
		nil,
		nil,
		nil,
	)

	router := setupRouter()
	router.GET("/crimes", controller.ListCrimes)

	params := usecases.ListCrimesParams{
		Page:  1,
		Limit: 10,
	}

	expectedCrimes := &entities.CrimeList{
		Items: []entities.Crime{
			{
				ID:          1,
				Title:       "Test Crime 1",
				Description: "Test Description 1",
				Type:        "ROBO",
				Status:      "ACTIVE",
				Location: entities.Location{
					Latitude:  -34.603722,
					Longitude: -58.381592,
					Address:   "Test Address 1",
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		Total: 1,
	}

	mockListUseCase.On("Execute", mock.Anything, params).Return(expectedCrimes, nil)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/crimes?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockListUseCase.AssertExpectations(t)
}

func TestCrimeController_UpdateCrimeStatus(t *testing.T) {
	// Arrange
	mockUpdateStatusUseCase := new(mocks.UpdateCrimeStatusUseCase)
	controller := controllers.NewCrimeController(
		nil,
		nil,
		mockUpdateStatusUseCase,
		nil,
		nil,
		nil,
	)

	router := setupRouter()
	router.PATCH("/crimes/:id/status", controller.UpdateCrimeStatus)

	input := usecases.UpdateCrimeStatusInput{
		UUID:   "123",
		Status: "INACTIVE",
	}

	mockUpdateStatusUseCase.On("Execute", mock.Anything, input).Return(nil)

	// Act
	reqBody := controllers.UpdateStatusRequest{
		Status: "INACTIVE",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/crimes/123/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "valid-key")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockUpdateStatusUseCase.AssertExpectations(t)
}

func TestCrimeController_GetCrimeStats(t *testing.T) {
	// Arrange
	mockGetStatsUseCase := new(mocks.GetCrimeStatsUseCase)
	controller := controllers.NewCrimeController(
		nil,
		nil,
		nil,
		nil,
		mockGetStatsUseCase,
		nil,
	)

	expectedStats := &entities.CrimeStats{
		TotalCrimes:    100,
		ActiveCrimes:   50,
		InactiveCrimes: 30,
		CrimesByType: map[string]int64{
			"robo":      60,
			"asalto":    30,
			"homicidio": 10,
		},
		CrimesByStatus: map[string]int64{
			"ACTIVE":   50,
			"INACTIVE": 30,
			"DELETED":  20,
		},
		CrimesByLocation: map[string]int64{},
		CrimesByAddress:  map[string]int64{},
		LastUpdate:       time.Now(),
	}

	mockGetStatsUseCase.On("Execute", mock.Anything).Return(expectedStats, nil)

	router := setupRouter()
	router.GET("/crimes/stats", controller.GetCrimeStats)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/crimes/stats", nil)
	req.Header.Set("X-API-Key", "valid-key")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockGetStatsUseCase.AssertExpectations(t)
}

func TestCrimeController_DeleteCrime(t *testing.T) {
	// Arrange
	mockDeleteUseCase := new(mocks.DeleteCrimeUseCase)
	controller := controllers.NewCrimeController(
		nil,
		nil,
		nil,
		mockDeleteUseCase,
		nil,
		nil,
	)

	mockDeleteUseCase.On("Execute", mock.Anything, "123").Return(nil)

	router := setupRouter()
	router.DELETE("/crimes/:id", controller.DeleteCrime)

	// Act
	req := httptest.NewRequest(http.MethodDelete, "/crimes/123", nil)
	req.Header.Set("X-API-Key", "valid-key")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	mockDeleteUseCase.AssertExpectations(t)
}

func TestCrimeController_GetCrime(t *testing.T) {
	// Arrange
	mockGetUseCase := new(mocks.GetCrimeUseCase)
	controller := controllers.NewCrimeController(
		nil,
		nil,
		nil,
		nil,
		nil,
		mockGetUseCase,
	)

	expectedCrime := &entities.Crime{
		ID:          1,
		Title:       "Test Crime",
		Description: "Test Description",
		Type:        "ROBO",
		Status:      "ACTIVE",
		Location: entities.Location{
			Latitude:  -34.603722,
			Longitude: -58.381592,
			Address:   "Test Address",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockGetUseCase.On("Execute", mock.Anything, "123").Return(expectedCrime, nil)

	router := setupRouter()
	router.GET("/crimes/:id", controller.GetCrime)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/crimes/123", nil)
	req.Header.Set("X-API-Key", "valid-key")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockGetUseCase.AssertExpectations(t)
}

func TestCrimeController_ListCrimes_Unauthorized(t *testing.T) {
	// Arrange
	mockCreateUseCase := new(mocks.CreateCrimeUseCase)
	mockListUseCase := new(mocks.ListCrimesUseCase)
	mockUpdateStatusUseCase := new(mocks.UpdateCrimeStatusUseCase)
	mockDeleteUseCase := new(mocks.DeleteCrimeUseCase)
	mockGetStatsUseCase := new(mocks.GetCrimeStatsUseCase)
	mockGetCrimeUseCase := new(mocks.GetCrimeUseCase)

	controller := controllers.NewCrimeController(
		mockCreateUseCase,
		mockListUseCase,
		mockUpdateStatusUseCase,
		mockDeleteUseCase,
		mockGetStatsUseCase,
		mockGetCrimeUseCase,
	)

	params := usecases.ListCrimesParams{
		Page:  1,
		Limit: 10,
	}

	mockListUseCase.On("Execute", mock.Anything, params).Return((*entities.CrimeList)(nil), errors.New("unauthorized"))

	router := setupRouter()
	router.GET("/crimes", controller.ListCrimes)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/crimes?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockListUseCase.AssertExpectations(t)
}

func TestCrimeController_CreateCrime_Unauthorized(t *testing.T) {
	// Arrange
	router := setupRouter()
	createUseCase := new(mocks.CreateCrimeUseCase)
	listUseCase := new(mocks.ListCrimesUseCase)
	updateStatusUseCase := new(mocks.UpdateCrimeStatusUseCase)
	deleteUseCase := new(mocks.DeleteCrimeUseCase)
	getStatsUseCase := new(mocks.GetCrimeStatsUseCase)
	getCrimeUseCase := new(mocks.GetCrimeUseCase)

	controller := controllers.NewCrimeController(
		createUseCase,
		listUseCase,
		updateStatusUseCase,
		deleteUseCase,
		getStatsUseCase,
		getCrimeUseCase,
	)

	addressNumber := "123"
	city := "Buenos Aires"
	province := "Buenos Aires"
	country := "Argentina"
	zipCode := "1000"

	reqBody := controllers.CreateCrimeRequest{
		Title:         "Test Crime",
		Description:   "Test Description",
		Type:          "ROBBERY",
		Latitude:      -34.603722,
		Longitude:     -58.381592,
		Address:       "Test Address",
		AddressNumber: addressNumber,
		City:          city,
		Province:      province,
		Country:       country,
		ZipCode:       zipCode,
	}

	createUseCase.On("Execute", mock.Anything, mock.AnythingOfType("usecases.CreateCrimeInput")).Return((*entities.Crime)(nil), errors.New("unauthorized"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/crimes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.POST("/crimes", controller.CreateCrime)
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	createUseCase.AssertExpectations(t)
}

func TestCrimeController_CreateCrime_InvalidAPIKey(t *testing.T) {
	// Arrange
	router := setupRouter()
	createUseCase := new(mocks.CreateCrimeUseCase)
	listUseCase := new(mocks.ListCrimesUseCase)
	updateStatusUseCase := new(mocks.UpdateCrimeStatusUseCase)
	deleteUseCase := new(mocks.DeleteCrimeUseCase)
	getStatsUseCase := new(mocks.GetCrimeStatsUseCase)
	getCrimeUseCase := new(mocks.GetCrimeUseCase)

	controller := controllers.NewCrimeController(
		createUseCase,
		listUseCase,
		updateStatusUseCase,
		deleteUseCase,
		getStatsUseCase,
		getCrimeUseCase,
	)

	addressNumber := "123"
	city := "Buenos Aires"
	province := "Buenos Aires"
	country := "Argentina"
	zipCode := "1000"

	reqBody := controllers.CreateCrimeRequest{
		Title:         "Test Crime",
		Description:   "Test Description",
		Type:          "ROBBERY",
		Latitude:      -34.603722,
		Longitude:     -58.381592,
		Address:       "Test Address",
		AddressNumber: addressNumber,
		City:          city,
		Province:      province,
		Country:       country,
		ZipCode:       zipCode,
	}

	createUseCase.On("Execute", mock.Anything, mock.AnythingOfType("usecases.CreateCrimeInput")).Return((*entities.Crime)(nil), errors.New("invalid api key"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/crimes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "invalid-key")
	w := httptest.NewRecorder()

	router.POST("/crimes", controller.CreateCrime)
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCrimeController_UpdateCrimeStatus_Unauthorized(t *testing.T) {
	// Arrange
	router := setupRouter()
	createUseCase := new(mocks.CreateCrimeUseCase)
	listUseCase := new(mocks.ListCrimesUseCase)
	updateStatusUseCase := new(mocks.UpdateCrimeStatusUseCase)
	deleteUseCase := new(mocks.DeleteCrimeUseCase)
	getStatsUseCase := new(mocks.GetCrimeStatsUseCase)
	getCrimeUseCase := new(mocks.GetCrimeUseCase)

	controller := controllers.NewCrimeController(
		createUseCase,
		listUseCase,
		updateStatusUseCase,
		deleteUseCase,
		getStatsUseCase,
		getCrimeUseCase,
	)

	input := usecases.UpdateCrimeStatusInput{
		UUID:   "123",
		Status: "INACTIVE",
	}

	updateStatusUseCase.On("Execute", mock.Anything, input).Return(errors.New("unauthorized"))

	reqBody := map[string]string{
		"status": "inactive",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PATCH", "/api/crimes/123/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.PATCH("/api/crimes/:id/status", controller.UpdateCrimeStatus)
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	updateStatusUseCase.AssertExpectations(t)
}

func TestCrimeController_DeleteCrime_Unauthorized(t *testing.T) {
	// Arrange
	router := setupRouter()
	createUseCase := new(mocks.CreateCrimeUseCase)
	listUseCase := new(mocks.ListCrimesUseCase)
	updateStatusUseCase := new(mocks.UpdateCrimeStatusUseCase)
	deleteUseCase := new(mocks.DeleteCrimeUseCase)
	getStatsUseCase := new(mocks.GetCrimeStatsUseCase)
	getCrimeUseCase := new(mocks.GetCrimeUseCase)

	controller := controllers.NewCrimeController(
		createUseCase,
		listUseCase,
		updateStatusUseCase,
		deleteUseCase,
		getStatsUseCase,
		getCrimeUseCase,
	)

	deleteUseCase.On("Execute", mock.Anything, "1").Return(errors.New("unauthorized"))

	req := httptest.NewRequest("DELETE", "/api/crimes/1", nil)
	w := httptest.NewRecorder()

	router.DELETE("/api/crimes/:id", controller.DeleteCrime)
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	deleteUseCase.AssertExpectations(t)
}

func TestCrimeController_DeleteCrime_NotFound(t *testing.T) {
	// Arrange
	router := setupRouter()
	createUseCase := new(mocks.CreateCrimeUseCase)
	listUseCase := new(mocks.ListCrimesUseCase)
	updateStatusUseCase := new(mocks.UpdateCrimeStatusUseCase)
	deleteUseCase := new(mocks.DeleteCrimeUseCase)
	getStatsUseCase := new(mocks.GetCrimeStatsUseCase)
	getCrimeUseCase := new(mocks.GetCrimeUseCase)

	controller := controllers.NewCrimeController(
		createUseCase,
		listUseCase,
		updateStatusUseCase,
		deleteUseCase,
		getStatsUseCase,
		getCrimeUseCase,
	)

	deleteUseCase.On("Execute", mock.Anything, "999").Return(errors.New("crime not found"))

	req := httptest.NewRequest("DELETE", "/api/crimes/999", nil)
	req.Header.Set("X-API-Key", "valid-key")
	w := httptest.NewRecorder()

	router.DELETE("/api/crimes/:id", controller.DeleteCrime)
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	deleteUseCase.AssertExpectations(t)
}

func TestCreateCrime_Success(t *testing.T) {
	// Arrange
	mockCreateUseCase := new(mocks.CreateCrimeUseCase)
	controller := controllers.NewCrimeController(mockCreateUseCase, nil, nil, nil, nil, nil)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/crimes", controller.CreateCrime)

	addressNumber := "123"
	city := "Buenos Aires"
	province := "CABA"
	country := "Argentina"
	zipCode := "1000"

	reqBody := controllers.CreateCrimeRequest{
		Title:         "Robo a mano armada",
		Description:   "Descripción del robo",
		Type:          "ROBO",
		Latitude:      -34.603722,
		Longitude:     -58.381592,
		Address:       "Calle Falsa",
		AddressNumber: addressNumber,
		City:          city,
		Province:      province,
		Country:       country,
		ZipCode:       zipCode,
	}

	input := usecases.CreateCrimeInput{
		Title:         reqBody.Title,
		Description:   reqBody.Description,
		Type:          reqBody.Type,
		Latitude:      reqBody.Latitude,
		Longitude:     reqBody.Longitude,
		Address:       reqBody.Address,
		AddressNumber: reqBody.AddressNumber,
		City:          reqBody.City,
		Province:      reqBody.Province,
		Country:       reqBody.Country,
		ZipCode:       reqBody.ZipCode,
	}

	expectedCrime := &entities.Crime{
		ID:          1,
		Title:       input.Title,
		Description: input.Description,
		Type:        input.Type,
		Location: entities.Location{
			Latitude:      input.Latitude,
			Longitude:     input.Longitude,
			Address:       input.Address,
			AddressNumber: input.AddressNumber,
			City:          input.City,
			Province:      input.Province,
			Country:       input.Country,
			ZipCode:       input.ZipCode,
		},
		Status:    "ACTIVE",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	mockCreateUseCase.On("Execute", mock.Anything, input).Return(expectedCrime, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/crimes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "valid-key")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	mockCreateUseCase.AssertExpectations(t)
}

func TestListCrimes_Success(t *testing.T) {
	// Arrange
	mockListUseCase := new(mocks.ListCrimesUseCase)
	controller := controllers.NewCrimeController(nil, mockListUseCase, nil, nil, nil, nil)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/crimes", controller.ListCrimes)

	params := usecases.ListCrimesParams{
		Page:  1,
		Limit: 10,
	}

	crimes := []entities.Crime{
		{
			ID:          1,
			Title:       "Robo a mano armada",
			Description: "Descripción del robo",
			Type:        "robo",
			Location: entities.Location{
				Latitude:  -34.603722,
				Longitude: -58.381592,
				Address:   "Calle Falsa 123",
			},
			Status:    "pendiente",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	expectedResponse := &entities.CrimeList{
		Items: crimes,
		Total: 1,
	}

	// Act
	mockListUseCase.On("Execute", mock.Anything, params).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/crimes?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockListUseCase.AssertExpectations(t)
}

func TestGetCrime_Error(t *testing.T) {
	// Arrange
	mockGetUseCase := new(mocks.GetCrimeUseCase)
	controller := controllers.NewCrimeController(
		nil,
		nil,
		nil,
		nil,
		nil,
		mockGetUseCase,
	)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/crimes/:id", controller.GetCrime)

	// Act
	mockGetUseCase.On("Execute", mock.Anything, "123").Return(nil, errors.New("error al obtener el delito"))

	req := httptest.NewRequest(http.MethodGet, "/crimes/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockGetUseCase.AssertExpectations(t)
}

func TestUpdateCrimeStatus_Success(t *testing.T) {
	// Arrange
	mockUseCase := new(mocks.UpdateCrimeStatusUseCase)
	controller := controllers.NewCrimeController(nil, nil, mockUseCase, nil, nil, nil)

	uuid := "test-uuid"
	input := usecases.UpdateCrimeStatusInput{
		UUID:   uuid,
		Status: "INACTIVE",
	}

	mockUseCase.On("Execute", mock.Anything, input).Return(nil)

	// Create request
	reqBody := `{"status": "INACTIVE"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/crimes/"+uuid+"/status", strings.NewReader(reqBody))

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: uuid}}

	// Act
	controller.UpdateCrimeStatus(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUpdateCrimeStatus_InvalidUUID(t *testing.T) {
	// Arrange
	mockUpdateStatusUseCase := new(mocks.UpdateCrimeStatusUseCase)
	controller := controllers.NewCrimeController(nil, nil, mockUpdateStatusUseCase, nil, nil, nil)

	// Si se llama, que retorne un error (no debería llamarse en este caso)
	mockUpdateStatusUseCase.On("Execute", mock.Anything, mock.Anything).Return(errors.New("should not be called"))

	// Create request
	reqBody := `{"status": "INACTIVE"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/crimes/status", strings.NewReader(reqBody))

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	controller.UpdateCrimeStatus(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteCrime_Success(t *testing.T) {
	// Arrange
	mockDeleteUseCase := new(mocks.DeleteCrimeUseCase)
	controller := controllers.NewCrimeController(
		nil,
		nil,
		nil,
		mockDeleteUseCase,
		nil,
		nil,
	)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/crimes/:id", controller.DeleteCrime)

	// Act
	mockDeleteUseCase.On("Execute", mock.Anything, "123").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/crimes/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	mockDeleteUseCase.AssertExpectations(t)
}

func TestGetCrimeStats_Success(t *testing.T) {
	// Arrange
	mockGetStatsUseCase := new(mocks.GetCrimeStatsUseCase)
	controller := controllers.NewCrimeController(
		nil,
		nil,
		nil,
		nil,
		mockGetStatsUseCase,
		nil,
	)

	router := setupRouter()
	router.GET("/crimes/stats", controller.GetCrimeStats)

	expectedStats := &entities.CrimeStats{
		TotalCrimes:      100,
		ActiveCrimes:     50,
		InactiveCrimes:   30,
		CrimesByType:     map[string]int64{"ROBO": 50, "HURTO": 30, "VIOLENTO": 20},
		CrimesByStatus:   map[string]int64{"ACTIVE": 70, "RESOLVED": 30},
		CrimesByLocation: map[string]int64{},
		CrimesByAddress:  map[string]int64{},
		LastUpdate:       time.Now(),
	}

	mockGetStatsUseCase.On("Execute", mock.Anything).Return(expectedStats, nil)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/crimes/stats", nil)
	req.Header.Set("X-API-Key", "valid-key")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockGetStatsUseCase.AssertExpectations(t)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, float64(100), response["total_crimes"])
	assert.NotNil(t, response["crimes_by_type"])
	assert.NotNil(t, response["crimes_by_status"])
	assert.NotNil(t, response["crimes_by_location"])
	assert.NotNil(t, response["crimes_by_address"])
}

func TestGetCrimeStats_Error(t *testing.T) {
	// Arrange
	mockGetStatsUseCase := new(mocks.GetCrimeStatsUseCase)
	controller := controllers.NewCrimeController(
		nil,
		nil,
		nil,
		nil,
		mockGetStatsUseCase,
		nil,
	)

	router := setupRouter()
	router.GET("/crimes/stats", controller.GetCrimeStats)

	expectedError := errors.New("error getting crime stats")
	mockGetStatsUseCase.On("Execute", mock.Anything).Return((*entities.CrimeStats)(nil), expectedError)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/crimes/stats", nil)
	req.Header.Set("X-API-Key", "valid-key")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockGetStatsUseCase.AssertExpectations(t)
}

// Casos de error
func TestCreateCrime_Error(t *testing.T) {
	// Arrange
	mockCreateUseCase := new(mocks.CreateCrimeUseCase)
	controller := controllers.NewCrimeController(mockCreateUseCase, nil, nil, nil, nil, nil)

	router := setupRouter()
	router.POST("/crimes", controller.CreateCrime)

	addressNumber := "123"
	city := "Buenos Aires"
	province := "Buenos Aires"
	country := "Argentina"
	zipCode := "1000"

	reqBody := controllers.CreateCrimeRequest{
		Title:         "Robo a mano armada",
		Description:   "Descripción del robo",
		Type:          "ROBO",
		Latitude:      -34.603722,
		Longitude:     -58.381592,
		Address:       "Calle Falsa 123",
		AddressNumber: addressNumber,
		City:          city,
		Province:      province,
		Country:       country,
		ZipCode:       zipCode,
	}

	// Act
	mockCreateUseCase.On("Execute", mock.Anything, mock.AnythingOfType("usecases.CreateCrimeInput")).Return(nil, errors.New("error al crear el delito"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/crimes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "valid-key")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockCreateUseCase.AssertExpectations(t)
}

func TestListCrimes_Error(t *testing.T) {
	// Arrange
	mockListUseCase := new(mocks.ListCrimesUseCase)
	controller := controllers.NewCrimeController(nil, mockListUseCase, nil, nil, nil, nil)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/crimes", controller.ListCrimes)

	params := usecases.ListCrimesParams{
		Page:  1,
		Limit: 10,
	}

	// Act
	mockListUseCase.On("Execute", mock.Anything, params).Return(nil, errors.New("error al listar los delitos"))

	req := httptest.NewRequest(http.MethodGet, "/crimes?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockListUseCase.AssertExpectations(t)
}

func TestUpdateCrimeStatus_Error(t *testing.T) {
	// Arrange
	mockUpdateStatusUseCase := new(mocks.UpdateCrimeStatusUseCase)
	controller := controllers.NewCrimeController(
		nil,
		nil,
		mockUpdateStatusUseCase,
		nil,
		nil,
		nil,
	)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/crimes/:id/status", controller.UpdateCrimeStatus)

	input := usecases.UpdateCrimeStatusInput{
		UUID:   "123",
		Status: "RESUELTO",
	}

	// Act
	mockUpdateStatusUseCase.On("Execute", mock.Anything, input).Return(errors.New("error al actualizar el estado"))

	body, _ := json.Marshal(map[string]string{"status": "RESUELTO"})
	req := httptest.NewRequest(http.MethodPut, "/crimes/123/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "valid-key")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUpdateStatusUseCase.AssertExpectations(t)
}

func TestDeleteCrime_Error(t *testing.T) {
	// Arrange
	mockDeleteUseCase := new(mocks.DeleteCrimeUseCase)
	controller := controllers.NewCrimeController(
		nil,
		nil,
		nil,
		mockDeleteUseCase,
		nil,
		nil,
	)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/crimes/:id", controller.DeleteCrime)

	// Act
	mockDeleteUseCase.On("Execute", mock.Anything, "123").Return(errors.New("error al eliminar el delito"))

	req := httptest.NewRequest(http.MethodDelete, "/crimes/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockDeleteUseCase.AssertExpectations(t)
}
