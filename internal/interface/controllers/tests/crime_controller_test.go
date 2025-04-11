package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/usecases"
	"go-crime_map_backend/internal/interface/controllers"
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

func TestCreateCrime(t *testing.T) {
	router := setupTestRouter()
	mockCreateUseCase := new(MockCreateCrimeUseCase)
	mockListUseCase := new(MockListCrimesUseCase)
	mockUpdateStatusUseCase := new(MockUpdateCrimeStatusUseCase)
	mockDeleteUseCase := new(MockDeleteCrimeUseCase)
	mockGetStatsUseCase := new(MockGetCrimeStatsUseCase)
	mockGetCrimeUseCase := new(MockGetCrimeUseCase)

	controller := controllers.NewCrimeController(
		mockCreateUseCase,
		mockListUseCase,
		mockUpdateStatusUseCase,
		mockDeleteUseCase,
		mockGetStatsUseCase,
		mockGetCrimeUseCase,
	)

	router.POST("/crimes", controller.CreateCrime)

	// Crear variables para los campos opcionales
	addressNumber := "1234"
	city := "Buenos Aires"
	province := "Buenos Aires"
	country := "Argentina"
	zipCode := "1000"

	tests := []struct {
		name           string
		input          controllers.CreateCrimeRequest
		mockSetup      func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "creación exitosa",
			input: controllers.CreateCrimeRequest{
				Title:         "Robo a mano armada",
				Description:   "Robo a mano armada en comercio",
				Type:          "ROBO",
				Latitude:      -34.603722,
				Longitude:     -58.381592,
				Address:       "Av. Corrientes",
				AddressNumber: utils.StringPtr(addressNumber),
				City:          utils.StringPtr(city),
				Province:      utils.StringPtr(province),
				Country:       utils.StringPtr(country),
				ZipCode:       utils.StringPtr(zipCode),
			},
			mockSetup: func() {
				mockCreateUseCase.On("Execute", mock.Anything, usecases.CreateCrimeInput{
					Title:         "Robo a mano armada",
					Description:   "Robo a mano armada en comercio",
					Type:          "ROBO",
					Latitude:      -34.603722,
					Longitude:     -58.381592,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr(addressNumber),
					City:          utils.StringPtr(city),
					Province:      utils.StringPtr(province),
					Country:       utils.StringPtr(country),
					ZipCode:       utils.StringPtr(zipCode),
				}).Return(
					&entities.Crime{
						ID:          "123",
						Title:       "Robo a mano armada",
						Description: "Robo a mano armada en comercio",
						Type:        "ROBO",
						Status:      "ACTIVE",
						Location: entities.Location{
							Latitude:      -34.603722,
							Longitude:     -58.381592,
							Address:       "Av. Corrientes",
							AddressNumber: utils.StringPtr(addressNumber),
							City:          utils.StringPtr(city),
							Province:      utils.StringPtr(province),
							Country:       utils.StringPtr(country),
							ZipCode:       utils.StringPtr(zipCode),
						},
					}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "error - datos inválidos",
			input: controllers.CreateCrimeRequest{
				Title:         "",
				Description:   "",
				Type:          "",
				Latitude:      0,
				Longitude:     0,
				Address:       "",
				AddressNumber: nil,
				City:          nil,
				Province:      nil,
				Country:       nil,
				ZipCode:       nil,
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Key: 'CreateCrimeRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/crimes", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				var response controllers.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.Error, tt.expectedError)
			}
		})
	}
}

func TestListCrimes(t *testing.T) {
	router := setupTestRouter()
	mockCreateUseCase := new(MockCreateCrimeUseCase)
	mockListUseCase := new(MockListCrimesUseCase)
	mockUpdateStatusUseCase := new(MockUpdateCrimeStatusUseCase)
	mockDeleteUseCase := new(MockDeleteCrimeUseCase)
	mockGetStatsUseCase := new(MockGetCrimeStatsUseCase)
	mockGetCrimeUseCase := new(MockGetCrimeUseCase)

	controller := controllers.NewCrimeController(
		mockCreateUseCase,
		mockListUseCase,
		mockUpdateStatusUseCase,
		mockDeleteUseCase,
		mockGetStatsUseCase,
		mockGetCrimeUseCase,
	)

	router.GET("/crimes", controller.ListCrimes)

	tests := []struct {
		name           string
		query          string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:  "listado exitoso",
			query: "?page=1&limit=10",
			mockSetup: func() {
				mockListUseCase.On("Execute", mock.Anything, usecases.ListCrimesParams{
					Page:  1,
					Limit: 10,
				}).Return(
					&entities.CrimeList{
						Items: []entities.Crime{
							{
								ID:          "123",
								Title:       "Robo a mano armada",
								Description: "Robo a mano armada en comercio",
								Type:        "ROBO",
								Status:      "ACTIVE",
								Location: entities.Location{
									Latitude:      -34.603722,
									Longitude:     -58.381592,
									Address:       "Av. Corrientes",
									AddressNumber: utils.StringPtr("1234"),
									City:          utils.StringPtr("Buenos Aires"),
									Province:      utils.StringPtr("Buenos Aires"),
									Country:       utils.StringPtr("Argentina"),
									ZipCode:       utils.StringPtr("1000"),
								},
							},
						},
						Total: 1,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			req := httptest.NewRequest("GET", "/crimes"+tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestUpdateCrimeStatus(t *testing.T) {
	router := setupTestRouter()
	mockCreateUseCase := new(MockCreateCrimeUseCase)
	mockListUseCase := new(MockListCrimesUseCase)
	mockUpdateStatusUseCase := new(MockUpdateCrimeStatusUseCase)
	mockDeleteUseCase := new(MockDeleteCrimeUseCase)
	mockGetStatsUseCase := new(MockGetCrimeStatsUseCase)
	mockGetCrimeUseCase := new(MockGetCrimeUseCase)

	controller := controllers.NewCrimeController(
		mockCreateUseCase,
		mockListUseCase,
		mockUpdateStatusUseCase,
		mockDeleteUseCase,
		mockGetStatsUseCase,
		mockGetCrimeUseCase,
	)

	router.PATCH("/crimes/:id/status", controller.UpdateCrimeStatus)

	tests := []struct {
		name           string
		crimeID        string
		input          controllers.UpdateStatusRequest
		mockSetup      func()
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "actualizar estado exitosamente",
			crimeID: "123",
			input: controllers.UpdateStatusRequest{
				Status: "INACTIVE",
			},
			mockSetup: func() {
				mockUpdateStatusUseCase.On("Execute", mock.Anything, usecases.UpdateCrimeStatusInput{
					ID:     "123",
					Status: "INACTIVE",
				}).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "error - delito no encontrado",
			crimeID: "456",
			input: controllers.UpdateStatusRequest{
				Status: "INACTIVE",
			},
			mockSetup: func() {
				mockUpdateStatusUseCase.On("Execute", mock.Anything, usecases.UpdateCrimeStatusInput{
					ID:     "456",
					Status: "INACTIVE",
				}).Return(fmt.Errorf("delito no encontrado"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "delito no encontrado",
		},
		{
			name:    "error - delito ya eliminado",
			crimeID: "789",
			input: controllers.UpdateStatusRequest{
				Status: "INACTIVE",
			},
			mockSetup: func() {
				mockUpdateStatusUseCase.On("Execute", mock.Anything, usecases.UpdateCrimeStatusInput{
					ID:     "789",
					Status: "INACTIVE",
				}).Return(fmt.Errorf("el delito ya ha sido eliminado"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "el delito ya ha sido eliminado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("PATCH", "/crimes/"+tt.crimeID+"/status", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				var response controllers.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.Error, tt.expectedError)
			}
		})
	}
}

func TestDeleteCrime(t *testing.T) {
	router := setupTestRouter()
	mockCreateUseCase := new(MockCreateCrimeUseCase)
	mockListUseCase := new(MockListCrimesUseCase)
	mockUpdateStatusUseCase := new(MockUpdateCrimeStatusUseCase)
	mockDeleteUseCase := new(MockDeleteCrimeUseCase)
	mockGetStatsUseCase := new(MockGetCrimeStatsUseCase)
	mockGetCrimeUseCase := new(MockGetCrimeUseCase)

	controller := controllers.NewCrimeController(
		mockCreateUseCase,
		mockListUseCase,
		mockUpdateStatusUseCase,
		mockDeleteUseCase,
		mockGetStatsUseCase,
		mockGetCrimeUseCase,
	)

	router.DELETE("/crimes/:id", controller.DeleteCrime)

	tests := []struct {
		name           string
		crimeID        string
		mockSetup      func()
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "eliminar delito exitosamente",
			crimeID: "123",
			mockSetup: func() {
				mockDeleteUseCase.On("Execute", mock.Anything, "123").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "error - delito no encontrado",
			crimeID: "456",
			mockSetup: func() {
				mockDeleteUseCase.On("Execute", mock.Anything, "456").Return(fmt.Errorf("delito no encontrado"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "delito no encontrado",
		},
		{
			name:    "error - delito ya eliminado",
			crimeID: "789",
			mockSetup: func() {
				mockDeleteUseCase.On("Execute", mock.Anything, "789").Return(fmt.Errorf("el delito ya ha sido eliminado"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "el delito ya ha sido eliminado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest("DELETE", "/crimes/"+tt.crimeID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				var response controllers.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.Error, tt.expectedError)
			}
		})
	}
}

func TestGetCrimeStats(t *testing.T) {
	router := setupTestRouter()
	mockCreateUseCase := new(MockCreateCrimeUseCase)
	mockListUseCase := new(MockListCrimesUseCase)
	mockUpdateStatusUseCase := new(MockUpdateCrimeStatusUseCase)
	mockDeleteUseCase := new(MockDeleteCrimeUseCase)
	mockGetStatsUseCase := new(MockGetCrimeStatsUseCase)
	mockGetCrimeUseCase := new(MockGetCrimeUseCase)

	controller := controllers.NewCrimeController(
		mockCreateUseCase,
		mockListUseCase,
		mockUpdateStatusUseCase,
		mockDeleteUseCase,
		mockGetStatsUseCase,
		mockGetCrimeUseCase,
	)

	router.GET("/crimes/stats", controller.GetCrimeStats)

	tests := []struct {
		name           string
		mockSetup      func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "obtener estadísticas exitosamente",
			mockSetup: func() {
				t.Log("Configurando mock para caso exitoso")
				mockGetStatsUseCase.On("Execute", mock.Anything).Return(
					&entities.CrimeStats{
						TotalCrimes:      10,
						ActiveCrimes:     5,
						InactiveCrimes:   3,
						CrimesByType:     map[string]int64{"ROBO": 5, "ASALTO": 3, "HURTO": 2},
						CrimesByStatus:   map[string]int64{"ACTIVO": 5, "INACTIVO": 3, "ELIMINADO": 2},
						CrimesByLocation: map[string]int64{"CABA": 5, "GBA": 3, "INTERIOR": 2},
						CrimesByAddress:  map[string]int64{"AV CORRIENTES": 5, "AV RIVADAVIA": 3, "AV CABILDO": 2},
						LastUpdate:       time.Now(),
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "error al obtener estadísticas",
			mockSetup: func() {
				t.Log("Configurando mock para caso de error")
				mockGetStatsUseCase.On("Execute", mock.Anything).Return(nil, fmt.Errorf("error al obtener estadísticas"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "error al obtener estadísticas",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Iniciando test case: %s", tt.name)
			// Limpiar las expectativas del mock antes de configurarlo
			mockGetStatsUseCase.ExpectedCalls = nil
			tt.mockSetup()

			req := httptest.NewRequest("GET", "/crimes/stats", nil)
			w := httptest.NewRecorder()
			t.Log("Ejecutando request HTTP")
			router.ServeHTTP(w, req)

			t.Logf("Verificando código de estado HTTP. Esperado: %d, Obtenido: %d", tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				t.Log("Verificando respuesta de error")
				var response controllers.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.Error, tt.expectedError)
			} else {
				t.Log("Verificando respuesta exitosa")
				var response entities.CrimeStats
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, int64(10), response.TotalCrimes)
				assert.Equal(t, int64(5), response.ActiveCrimes)
				assert.Equal(t, int64(3), response.InactiveCrimes)
				assert.Equal(t, map[string]int64{"ROBO": 5, "ASALTO": 3, "HURTO": 2}, response.CrimesByType)
				assert.Equal(t, map[string]int64{"ACTIVO": 5, "INACTIVO": 3, "ELIMINADO": 2}, response.CrimesByStatus)
				assert.Equal(t, map[string]int64{"CABA": 5, "GBA": 3, "INTERIOR": 2}, response.CrimesByLocation)
				assert.Equal(t, map[string]int64{"AV CORRIENTES": 5, "AV RIVADAVIA": 3, "AV CABILDO": 2}, response.CrimesByAddress)
				assert.NotZero(t, response.LastUpdate)
			}
			t.Logf("Test case %s completado", tt.name)
		})
	}
}
