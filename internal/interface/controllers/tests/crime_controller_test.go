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
	"go-crime_map_backend/internal/interface/controllers"
	"go-crime_map_backend/internal/mocks"
	"go-crime_map_backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCreateCrimeUseCase es un mock para el caso de uso de creación de delitos
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

// MockListCrimesUseCase es un mock para el caso de uso de listado de delitos
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

func (m *MockDeleteCrimeUseCase) Execute(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockGetCrimeStatsUseCase es un mock para el caso de uso de obtención de estadísticas
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

// MockGetCrimeUseCase es un mock para el caso de uso de obtención de un delito
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

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestCrimeController_CreateCrime(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Crear delito exitosamente", func(t *testing.T) {
		mockCreateUseCase := new(mocks.MockCreateCrimeUseCase)
		controller := controllers.NewCrimeController(
			mockCreateUseCase,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		req := controllers.CreateCrimeRequest{
			Title:         "Robo",
			Description:   "Robo en tienda",
			Type:          "ROBBERY",
			Latitude:      40.7128,
			Longitude:     -74.0060,
			Address:       "Av. Corrientes",
			AddressNumber: utils.StringPtr("1234"),
			City:          utils.StringPtr("Buenos Aires"),
			Province:      utils.StringPtr("Buenos Aires"),
			Country:       utils.StringPtr("Argentina"),
			ZipCode:       utils.StringPtr("1000"),
		}

		expectedCrime := &entities.Crime{
			ID:          "1",
			Title:       req.Title,
			Description: req.Description,
			Type:        req.Type,
			Status:      "ACTIVE",
			Location: entities.Location{
				Latitude:      req.Latitude,
				Longitude:     req.Longitude,
				Address:       req.Address,
				AddressNumber: req.AddressNumber,
				City:          req.City,
				Province:      req.Province,
				Country:       req.Country,
				ZipCode:       req.ZipCode,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockCreateUseCase.On("Execute", mock.Anything, mock.Anything).Return(expectedCrime, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/crimes", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.CreateCrime(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response entities.Crime
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedCrime.ID, response.ID)
		assert.Equal(t, expectedCrime.Title, response.Title)
		assert.Equal(t, expectedCrime.Description, response.Description)
		assert.Equal(t, expectedCrime.Type, response.Type)
		assert.Equal(t, expectedCrime.Status, response.Status)
		assert.Equal(t, expectedCrime.Location.Latitude, response.Location.Latitude)
		assert.Equal(t, expectedCrime.Location.Longitude, response.Location.Longitude)
		assert.Equal(t, expectedCrime.Location.Address, response.Location.Address)
		assert.Equal(t, expectedCrime.Location.AddressNumber, response.Location.AddressNumber)

		mockCreateUseCase.AssertExpectations(t)
	})

	t.Run("Error al crear delito - Datos inválidos", func(t *testing.T) {
		mockCreateUseCase := new(mocks.MockCreateCrimeUseCase)
		controller := controllers.NewCrimeController(
			mockCreateUseCase,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		req := controllers.CreateCrimeRequest{
			// Datos incompletos
			Title: "Robo",
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/crimes", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.CreateCrime(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockCreateUseCase.AssertNotCalled(t, "Execute")
	})

	t.Run("Error al crear delito - Error interno", func(t *testing.T) {
		mockCreateUseCase := new(mocks.MockCreateCrimeUseCase)
		controller := controllers.NewCrimeController(
			mockCreateUseCase,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		req := controllers.CreateCrimeRequest{
			Title:         "Robo",
			Description:   "Robo en tienda",
			Type:          "ROBBERY",
			Latitude:      40.7128,
			Longitude:     -74.0060,
			Address:       "Av. Corrientes",
			AddressNumber: utils.StringPtr("1234"),
			City:          utils.StringPtr("Buenos Aires"),
			Province:      utils.StringPtr("Buenos Aires"),
			Country:       utils.StringPtr("Argentina"),
			ZipCode:       utils.StringPtr("1000"),
		}

		expectedError := errors.New("error interno")
		mockCreateUseCase.On("Execute", mock.Anything, mock.Anything).Return(nil, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPost, "/crimes", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.CreateCrime(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response controllers.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedError.Error(), response.Error)

		mockCreateUseCase.AssertExpectations(t)
	})
}

func TestCrimeController_ListCrimes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Listar delitos exitosamente", func(t *testing.T) {
		mockListUseCase := new(mocks.MockListCrimesUseCase)
		controller := controllers.NewCrimeController(
			nil,
			mockListUseCase,
			nil,
			nil,
			nil,
			nil,
		)

		expectedCrimes := &entities.CrimeList{
			Items: []entities.Crime{
				{
					ID:          "1",
					Title:       "Robo",
					Description: "Robo en tienda",
					Type:        "ROBBERY",
					Status:      "ACTIVE",
					Location: entities.Location{
						Latitude:      40.7128,
						Longitude:     -74.0060,
						Address:       "Av. Corrientes",
						AddressNumber: utils.StringPtr("1234"),
					},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			Total: 1,
		}

		mockListUseCase.On("Execute", mock.Anything, mock.Anything).Return(expectedCrimes, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/crimes?page=1&limit=10", nil)

		controller.ListCrimes(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entities.CrimeList
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedCrimes.Total, response.Total)
		assert.Equal(t, len(expectedCrimes.Items), len(response.Items))
		assert.Equal(t, expectedCrimes.Items[0].ID, response.Items[0].ID)

		mockListUseCase.AssertExpectations(t)
	})

	t.Run("Error al listar delitos - Parámetros inválidos", func(t *testing.T) {
		mockListUseCase := new(mocks.MockListCrimesUseCase)
		controller := controllers.NewCrimeController(
			nil,
			mockListUseCase,
			nil,
			nil,
			nil,
			nil,
		)

		expectedCrimes := &entities.CrimeList{
			Items: []entities.Crime{},
			Total: 0,
		}

		mockListUseCase.On("Execute", mock.Anything, mock.Anything).Return(expectedCrimes, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/crimes?page=invalid&limit=invalid", nil)

		controller.ListCrimes(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockListUseCase.AssertExpectations(t)
	})

	t.Run("Error al listar delitos - Error interno", func(t *testing.T) {
		mockListUseCase := new(mocks.MockListCrimesUseCase)
		controller := controllers.NewCrimeController(
			nil,
			mockListUseCase,
			nil,
			nil,
			nil,
			nil,
		)

		expectedError := errors.New("error al listar delitos")
		mockListUseCase.On("Execute", mock.Anything, mock.Anything).Return(nil, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/crimes?page=1&limit=10", nil)

		controller.ListCrimes(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response controllers.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedError.Error(), response.Error)

		mockListUseCase.AssertExpectations(t)
	})

	t.Run("Listar delitos exitosamente con todos los parámetros", func(t *testing.T) {
		mockListUseCase := new(mocks.MockListCrimesUseCase)
		controller := controllers.NewCrimeController(
			nil,
			mockListUseCase,
			nil,
			nil,
			nil,
			nil,
		)

		expectedCrimes := &entities.CrimeList{
			Items: []entities.Crime{
				{
					ID:          "1",
					Title:       "Robo",
					Description: "Robo en tienda",
					Type:        "ROBBERY",
					Status:      "ACTIVE",
				},
			},
			Total: 1,
		}

		mockListUseCase.On("Execute", mock.Anything, mock.Anything).Return(expectedCrimes, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/crimes?page=1&limit=10&type=ROBBERY&status=ACTIVE&start_date=2024-01-01&end_date=2024-12-31", nil)

		controller.ListCrimes(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response entities.CrimeList
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedCrimes.Total, response.Total)
		assert.Equal(t, len(expectedCrimes.Items), len(response.Items))

		mockListUseCase.AssertExpectations(t)
	})

	t.Run("Listar delitos - Parámetros de fecha inválidos", func(t *testing.T) {
		mockListUseCase := new(mocks.MockListCrimesUseCase)
		controller := controllers.NewCrimeController(
			nil,
			mockListUseCase,
			nil,
			nil,
			nil,
			nil,
		)

		expectedError := errors.New("fecha inválida")
		mockListUseCase.On("Execute", mock.Anything, mock.Anything).Return(nil, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/crimes?start_date=invalid&end_date=invalid", nil)

		controller.ListCrimes(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response controllers.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedError.Error(), response.Error)

		mockListUseCase.AssertExpectations(t)
	})

	t.Run("Listar delitos - Tipo de delito inválido", func(t *testing.T) {
		mockListUseCase := new(mocks.MockListCrimesUseCase)
		controller := controllers.NewCrimeController(
			nil,
			mockListUseCase,
			nil,
			nil,
			nil,
			nil,
		)

		expectedError := errors.New("tipo de delito inválido")
		mockListUseCase.On("Execute", mock.Anything, mock.Anything).Return(nil, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/crimes?type=INVALID", nil)

		controller.ListCrimes(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response controllers.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedError.Error(), response.Error)

		mockListUseCase.AssertExpectations(t)
	})

	t.Run("Error al listar delitos - fecha inicial inválida", func(t *testing.T) {
		mockUseCase := new(mocks.MockListCrimesUseCase)
		controller := controllers.NewCrimeController(nil, mockUseCase, nil, nil, nil, nil)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/?start_date=fecha-invalida", nil)

		expectedError := errors.New("fecha inicial inválida")
		mockUseCase.On("Execute", mock.Anything, mock.AnythingOfType("usecases.ListCrimesParams")).Return(nil, expectedError)

		controller.ListCrimes(c)

		assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error al listar delitos - fecha final inválida", func(t *testing.T) {
		mockUseCase := new(mocks.MockListCrimesUseCase)
		controller := controllers.NewCrimeController(nil, mockUseCase, nil, nil, nil, nil)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/?end_date=fecha-invalida", nil)

		expectedError := errors.New("fecha final inválida")
		mockUseCase.On("Execute", mock.Anything, mock.AnythingOfType("usecases.ListCrimesParams")).Return(nil, expectedError)

		controller.ListCrimes(c)

		assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error al listar delitos - radio inválido", func(t *testing.T) {
		mockUseCase := new(mocks.MockListCrimesUseCase)
		controller := controllers.NewCrimeController(nil, mockUseCase, nil, nil, nil, nil)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/?radius_km=no-numerico", nil)

		expectedError := errors.New("radio inválido")
		mockUseCase.On("Execute", mock.Anything, mock.AnythingOfType("usecases.ListCrimesParams")).Return(nil, expectedError)

		controller.ListCrimes(c)

		assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error al listar delitos - latitud inválida", func(t *testing.T) {
		mockUseCase := new(mocks.MockListCrimesUseCase)
		controller := controllers.NewCrimeController(nil, mockUseCase, nil, nil, nil, nil)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/?latitude=no-numerico", nil)

		expectedError := errors.New("latitud inválida")
		mockUseCase.On("Execute", mock.Anything, mock.AnythingOfType("usecases.ListCrimesParams")).Return(nil, expectedError)

		controller.ListCrimes(c)

		assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error al listar delitos - longitud inválida", func(t *testing.T) {
		mockUseCase := new(mocks.MockListCrimesUseCase)
		controller := controllers.NewCrimeController(nil, mockUseCase, nil, nil, nil, nil)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/?longitude=no-numerico", nil)

		expectedError := errors.New("longitud inválida")
		mockUseCase.On("Execute", mock.Anything, mock.AnythingOfType("usecases.ListCrimesParams")).Return(nil, expectedError)

		controller.ListCrimes(c)

		assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error al listar delitos - página inválida", func(t *testing.T) {
		mockUseCase := new(mocks.MockListCrimesUseCase)
		controller := controllers.NewCrimeController(nil, mockUseCase, nil, nil, nil, nil)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/?page=no-numerico", nil)

		expectedError := errors.New("página inválida")
		mockUseCase.On("Execute", mock.Anything, mock.AnythingOfType("usecases.ListCrimesParams")).Return(nil, expectedError)

		controller.ListCrimes(c)

		assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error al listar delitos - límite inválido", func(t *testing.T) {
		mockUseCase := new(mocks.MockListCrimesUseCase)
		controller := controllers.NewCrimeController(nil, mockUseCase, nil, nil, nil, nil)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/?limit=no-numerico", nil)

		expectedError := errors.New("límite inválido")
		mockUseCase.On("Execute", mock.Anything, mock.AnythingOfType("usecases.ListCrimesParams")).Return(nil, expectedError)

		controller.ListCrimes(c)

		assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
		mockUseCase.AssertExpectations(t)
	})
}

func TestCrimeController_UpdateCrimeStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Actualizar estado exitosamente", func(t *testing.T) {
		mockUpdateStatusUseCase := new(mocks.MockUpdateCrimeStatusUseCase)
		controller := controllers.NewCrimeController(
			nil,
			nil,
			mockUpdateStatusUseCase,
			nil,
			nil,
			nil,
		)

		req := controllers.UpdateStatusRequest{
			Status: "INACTIVE",
		}

		mockUpdateStatusUseCase.On("Execute", mock.Anything, mock.Anything).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		body, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPatch, "/crimes/1/status", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.UpdateCrimeStatus(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUpdateStatusUseCase.AssertExpectations(t)
	})

	t.Run("Error al actualizar estado - ID inválido", func(t *testing.T) {
		mockUpdateStatusUseCase := new(mocks.MockUpdateCrimeStatusUseCase)
		controller := controllers.NewCrimeController(
			nil,
			nil,
			mockUpdateStatusUseCase,
			nil,
			nil,
			nil,
		)

		req := controllers.UpdateStatusRequest{
			Status: "INACTIVE",
		}

		expectedError := errors.New("ID inválido")
		mockUpdateStatusUseCase.On("Execute", mock.Anything, mock.Anything).Return(expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: ""}}

		body, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPatch, "/crimes//status", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.UpdateCrimeStatus(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response controllers.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedError.Error(), response.Error)

		mockUpdateStatusUseCase.AssertExpectations(t)
	})

	t.Run("Error al actualizar estado - Error interno", func(t *testing.T) {
		mockUpdateStatusUseCase := new(mocks.MockUpdateCrimeStatusUseCase)
		controller := controllers.NewCrimeController(
			nil,
			nil,
			mockUpdateStatusUseCase,
			nil,
			nil,
			nil,
		)

		req := controllers.UpdateStatusRequest{
			Status: "INACTIVE",
		}

		expectedError := errors.New("error al actualizar estado")
		mockUpdateStatusUseCase.On("Execute", mock.Anything, mock.Anything).Return(expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		body, _ := json.Marshal(req)
		c.Request = httptest.NewRequest(http.MethodPatch, "/crimes/1/status", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.UpdateCrimeStatus(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response controllers.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedError.Error(), response.Error)

		mockUpdateStatusUseCase.AssertExpectations(t)
	})

	t.Run("Error al actualizar estado - JSON inválido", func(t *testing.T) {
		mockUseCase := new(mocks.MockUpdateCrimeStatusUseCase)
		controller := controllers.NewCrimeController(nil, nil, mockUseCase, nil, nil, nil)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodPut, "/", strings.NewReader("{json_invalido}"))

		controller.UpdateCrimeStatus(c)

		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	})

	t.Run("Error al actualizar estado - JSON vacío", func(t *testing.T) {
		mockUseCase := new(mocks.MockUpdateCrimeStatusUseCase)
		controller := controllers.NewCrimeController(nil, nil, mockUseCase, nil, nil, nil)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(""))

		controller.UpdateCrimeStatus(c)

		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	})

	t.Run("Error al actualizar estado - Campo status faltante", func(t *testing.T) {
		mockUseCase := new(mocks.MockUpdateCrimeStatusUseCase)
		controller := controllers.NewCrimeController(nil, nil, mockUseCase, nil, nil, nil)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"id": "123"}`))

		controller.UpdateCrimeStatus(c)

		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	})

	t.Run("Error al actualizar estado - Campo id faltante", func(t *testing.T) {
		mockUseCase := new(mocks.MockUpdateCrimeStatusUseCase)
		controller := controllers.NewCrimeController(nil, nil, mockUseCase, nil, nil, nil)

		// Configurar el mock para devolver un error cuando se llama con ID vacío
		mockUseCase.On("Execute", mock.Anything, usecases.UpdateCrimeStatusInput{
			ID:     "",
			Status: "INACTIVE",
		}).Return(errors.New("ID inválido"))

		// Crear el request con JSON sin el campo ID
		requestBody := `{"status": "INACTIVE"}`
		req, _ := http.NewRequest("PUT", "/crimes/status", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		controller.UpdateCrimeStatus(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Eliminar delito exitosamente", func(t *testing.T) {
		mockDeleteUseCase := new(mocks.MockDeleteCrimeUseCase)
		controller := controllers.NewCrimeController(
			nil,
			nil,
			nil,
			mockDeleteUseCase,
			nil,
			nil,
		)

		mockDeleteUseCase.On("Execute", mock.Anything, "1").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		c.Request = httptest.NewRequest(http.MethodDelete, "/crimes/1", nil)

		controller.DeleteCrime(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockDeleteUseCase.AssertExpectations(t)
	})

	t.Run("Error al eliminar delito - ID inválido", func(t *testing.T) {
		mockUseCase := new(mocks.MockDeleteCrimeUseCase)
		controller := controllers.NewCrimeController(nil, nil, nil, mockUseCase, nil, nil)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)

		expectedError := errors.New("ID inválido")
		mockUseCase.On("Execute", mock.Anything, "").Return(expectedError)

		controller.DeleteCrime(c)

		assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error al eliminar delito - Error interno", func(t *testing.T) {
		mockUseCase := new(mocks.MockDeleteCrimeUseCase)
		controller := controllers.NewCrimeController(nil, nil, nil, mockUseCase, nil, nil)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Params = []gin.Param{{Key: "id", Value: "123"}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)

		expectedError := errors.New("error interno")
		mockUseCase.On("Execute", mock.Anything, "123").Return(expectedError)

		controller.DeleteCrime(c)

		assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
		mockUseCase.AssertExpectations(t)
	})
}

func TestCrimeController_GetCrimeStats(t *testing.T) {
	t.Run("Error al obtener estadísticas - Error interno", func(t *testing.T) {
		mockUseCase := new(mocks.MockGetCrimeStatsUseCase)
		controller := controllers.NewCrimeController(nil, nil, nil, nil, mockUseCase, nil)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/stats", nil)

		expectedError := errors.New("error interno")
		mockUseCase.On("Execute", mock.Anything).Return(nil, expectedError)

		controller.GetCrimeStats(c)

		assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error al obtener estadísticas - Sin datos", func(t *testing.T) {
		mockUseCase := new(mocks.MockGetCrimeStatsUseCase)
		controller := controllers.NewCrimeController(nil, nil, nil, nil, mockUseCase, nil)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/stats", nil)

		mockUseCase.On("Execute", mock.Anything).Return(nil, nil)

		controller.GetCrimeStats(c)

		assert.Equal(t, http.StatusNotFound, c.Writer.Status())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Obtener estadísticas exitosamente", func(t *testing.T) {
		mockGetStatsUseCase := new(mocks.MockGetCrimeStatsUseCase)
		controller := controllers.NewCrimeController(
			nil,
			nil,
			nil,
			nil,
			mockGetStatsUseCase,
			nil,
		)

		expectedStats := &entities.CrimeStats{
			TotalCrimes:      10,
			ActiveCrimes:     5,
			InactiveCrimes:   3,
			CrimesByType:     map[string]int64{"ROBO": 5, "ASALTO": 3},
			CrimesByStatus:   map[string]int64{"ACTIVE": 5, "INACTIVE": 3},
			CrimesByLocation: map[string]int64{"CENTRO": 5, "SUR": 3},
			CrimesByAddress:  map[string]int64{"AV CORRIENTES": 5, "AV RIVADAVIA": 3},
			LastUpdate:       time.Now(),
		}

		mockGetStatsUseCase.On("Execute", mock.Anything).Return(expectedStats, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/crimes/stats", nil)

		controller.GetCrimeStats(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entities.CrimeStats
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedStats.TotalCrimes, response.TotalCrimes)
		assert.Equal(t, expectedStats.ActiveCrimes, response.ActiveCrimes)
		assert.Equal(t, expectedStats.InactiveCrimes, response.InactiveCrimes)
		assert.Equal(t, expectedStats.CrimesByType, response.CrimesByType)
		assert.Equal(t, expectedStats.CrimesByStatus, response.CrimesByStatus)
		assert.Equal(t, expectedStats.CrimesByLocation, response.CrimesByLocation)
		assert.Equal(t, expectedStats.CrimesByAddress, response.CrimesByAddress)

		mockGetStatsUseCase.AssertExpectations(t)
	})

	t.Run("Error al obtener estadísticas - Error de base de datos", func(t *testing.T) {
		mockGetStatsUseCase := new(mocks.MockGetCrimeStatsUseCase)
		controller := controllers.NewCrimeController(
			nil,
			nil,
			nil,
			nil,
			mockGetStatsUseCase,
			nil,
		)

		expectedError := errors.New("error de base de datos")
		mockGetStatsUseCase.On("Execute", mock.Anything).Return(nil, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/crimes/stats", nil)

		controller.GetCrimeStats(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response controllers.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedError.Error(), response.Error)

		mockGetStatsUseCase.AssertExpectations(t)
	})

	t.Run("Error al obtener estadísticas - Error de procesamiento", func(t *testing.T) {
		mockGetStatsUseCase := new(mocks.MockGetCrimeStatsUseCase)
		controller := controllers.NewCrimeController(
			nil,
			nil,
			nil,
			nil,
			mockGetStatsUseCase,
			nil,
		)

		expectedError := errors.New("error al procesar estadísticas")
		mockGetStatsUseCase.On("Execute", mock.Anything).Return(nil, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/crimes/stats", nil)

		controller.GetCrimeStats(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response controllers.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedError.Error(), response.Error)

		mockGetStatsUseCase.AssertExpectations(t)
	})
}

func TestCrimeController_GetCrime(t *testing.T) {
	t.Run("Error al obtener delito - ID inválido", func(t *testing.T) {
		mockUseCase := new(mocks.MockGetCrimeUseCase)
		controller := controllers.NewCrimeController(nil, nil, nil, nil, nil, mockUseCase)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		expectedError := errors.New("ID inválido")
		mockUseCase.On("Execute", mock.Anything, "").Return(nil, expectedError)

		controller.GetCrime(c)

		assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error al obtener delito - No encontrado", func(t *testing.T) {
		mockUseCase := new(mocks.MockGetCrimeUseCase)
		controller := controllers.NewCrimeController(nil, nil, nil, nil, nil, mockUseCase)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Params = []gin.Param{{Key: "id", Value: "123"}}
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		mockUseCase.On("Execute", mock.Anything, "123").Return(nil, nil)

		controller.GetCrime(c)

		assert.Equal(t, http.StatusNotFound, c.Writer.Status())
		mockUseCase.AssertExpectations(t)
	})
}

func TestParseFloat64(t *testing.T) {
	t.Run("Cadena vacía", func(t *testing.T) {
		result := controllers.ParseFloat64("")
		assert.Equal(t, 0.0, result)
	})

	t.Run("Número válido", func(t *testing.T) {
		result := controllers.ParseFloat64("3.14")
		assert.Equal(t, 3.14, result)
	})

	t.Run("Número entero", func(t *testing.T) {
		result := controllers.ParseFloat64("42")
		assert.Equal(t, 42.0, result)
	})

	t.Run("Número negativo", func(t *testing.T) {
		result := controllers.ParseFloat64("-3.14")
		assert.Equal(t, -3.14, result)
	})

	t.Run("Cadena inválida", func(t *testing.T) {
		result := controllers.ParseFloat64("no es un número")
		assert.Equal(t, 0.0, result)
	})

	t.Run("Cadena con espacios", func(t *testing.T) {
		result := controllers.ParseFloat64(" 3.14 ")
		assert.Equal(t, 3.14, result)
	})
}
