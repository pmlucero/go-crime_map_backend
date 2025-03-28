package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/usecases"
	"go-crime_map_backend/internal/interface/controllers"
	"go-crime_map_backend/internal/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (m *MockListCrimesUseCase) Execute(ctx context.Context, input usecases.ListCrimesInput) ([]entities.Crime, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Crime), args.Error(1)
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

func (m *MockGetCrimeStatsUseCase) Execute(ctx context.Context, input usecases.GetCrimeStatsInput) (*entities.CrimeStats, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CrimeStats), args.Error(1)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestCreateCrime(t *testing.T) {
	router := setupTestRouter()

	mockCreateUseCase := new(mocks.MockCreateCrimeUseCase)
	mockListUseCase := new(mocks.MockListCrimesUseCase)
	mockUpdateStatusUseCase := new(mocks.MockUpdateCrimeStatusUseCase)
	mockDeleteUseCase := new(mocks.MockDeleteCrimeUseCase)
	mockGetStatsUseCase := new(mocks.MockGetCrimeStatsUseCase)
	mockGetCrimeUseCase := new(mocks.MockGetCrimeUseCase)

	controller := controllers.NewCrimeController(
		mockCreateUseCase,
		mockListUseCase,
		mockUpdateStatusUseCase,
		mockDeleteUseCase,
		mockGetStatsUseCase,
		mockGetCrimeUseCase,
	)

	router.POST("/crimes", controller.CreateCrime)

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
				Title:       "Robo a mano armada",
				Description: "Robo a mano armada en comercio",
				Type:        "ROBO",
				Latitude:    -34.603722,
				Longitude:   -58.381592,
				Address:     "Av. Corrientes 1234",
			},
			mockSetup: func() {
				mockCreateUseCase.On("Execute", mock.Anything, usecases.CreateCrimeInput{
					Title:       "Robo a mano armada",
					Description: "Robo a mano armada en comercio",
					Type:        "ROBO",
					Latitude:    -34.603722,
					Longitude:   -58.381592,
					Address:     "Av. Corrientes 1234",
				}).Return(
					&entities.Crime{
						ID:          "123",
						Title:       "Robo a mano armada",
						Description: "Robo a mano armada en comercio",
						Type:        "ROBO",
						Status:      "ACTIVE",
						Location: entities.Location{
							Latitude:  -34.603722,
							Longitude: -58.381592,
							Address:   "Av. Corrientes 1234",
						},
					}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "error - datos inválidos",
			input: controllers.CreateCrimeRequest{
				Title:       "",
				Description: "",
				Type:        "",
				Latitude:    0,
				Longitude:   0,
				Address:     "",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "error de validación",
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

	mockCreateUseCase := new(mocks.MockCreateCrimeUseCase)
	mockListUseCase := new(mocks.MockListCrimesUseCase)
	mockUpdateStatusUseCase := new(mocks.MockUpdateCrimeStatusUseCase)
	mockDeleteUseCase := new(mocks.MockDeleteCrimeUseCase)
	mockGetStatsUseCase := new(mocks.MockGetCrimeStatsUseCase)
	mockGetCrimeUseCase := new(mocks.MockGetCrimeUseCase)

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
									Latitude:  -34.603722,
									Longitude: -58.381592,
									Address:   "Av. Corrientes 1234",
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

	controller := controllers.NewCrimeController(
		mockCreateUseCase,
		mockListUseCase,
		mockUpdateStatusUseCase,
		mockDeleteUseCase,
		mockGetStatsUseCase,
	)

	router.PATCH("/crimes/:id/status", controller.UpdateCrimeStatus)

	tests := []struct {
		name           string
		id             string
		payload        controllers.UpdateStatusRequest
		mockSetup      func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "actualizar estado exitosamente",
			id:   uuid.New().String(),
			payload: controllers.UpdateStatusRequest{
				Status: string(entities.CrimeStatusInactive),
			},
			mockSetup: func() {
				mockUpdateStatusUseCase.On("Execute", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "error - ID no proporcionado",
			id:   "",
			payload: controllers.UpdateStatusRequest{
				Status: string(entities.CrimeStatusInactive),
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "ID no proporcionado",
		},
		{
			name: "error - estado inválido",
			id:   uuid.New().String(),
			payload: controllers.UpdateStatusRequest{
				Status: "INVALID_STATUS",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "estado inválido",
		},
		{
			name: "error - delito no encontrado",
			id:   uuid.New().String(),
			payload: controllers.UpdateStatusRequest{
				Status: string(entities.CrimeStatusInactive),
			},
			mockSetup: func() {
				mockUpdateStatusUseCase.On("Execute", mock.Anything, mock.Anything).Return(usecases.ErrCrimeNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "el delito no existe",
		},
		{
			name: "error - delito ya eliminado",
			id:   uuid.New().String(),
			payload: controllers.UpdateStatusRequest{
				Status: string(entities.CrimeStatusInactive),
			},
			mockSetup: func() {
				mockUpdateStatusUseCase.On("Execute", mock.Anything, mock.Anything).Return(usecases.ErrCrimeAlreadyDeleted)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "el delito ya fue eliminado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Limpiar mocks antes de cada test
			mockUpdateStatusUseCase.ExpectedCalls = nil
			mockUpdateStatusUseCase.Calls = nil

			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("PATCH", "/crimes/"+tt.id+"/status", bytes.NewBuffer(payload))
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

	controller := controllers.NewCrimeController(
		mockCreateUseCase,
		mockListUseCase,
		mockUpdateStatusUseCase,
		mockDeleteUseCase,
		mockGetStatsUseCase,
	)

	crimes := router.Group("/crimes")
	{
		crimes.DELETE("", controller.DeleteCrime)
		crimes.DELETE("/:id", controller.DeleteCrime)
	}

	tests := []struct {
		name           string
		path           string
		mockSetup      func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "eliminar delito exitosamente",
			path: "/crimes/" + uuid.New().String(),
			mockSetup: func() {
				mockDeleteUseCase.On("Execute", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error - ID no proporcionado",
			path:           "/crimes",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "ID no proporcionado",
		},
		{
			name: "error - delito no encontrado",
			path: "/crimes/" + uuid.New().String(),
			mockSetup: func() {
				mockDeleteUseCase.On("Execute", mock.Anything, mock.Anything).Return(usecases.ErrCrimeNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "el delito no existe",
		},
		{
			name: "error - delito ya eliminado",
			path: "/crimes/" + uuid.New().String(),
			mockSetup: func() {
				mockDeleteUseCase.On("Execute", mock.Anything, mock.Anything).Return(usecases.ErrCrimeAlreadyDeleted)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "el delito ya fue eliminado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Limpiar mocks antes de cada test
			mockDeleteUseCase.ExpectedCalls = nil
			mockDeleteUseCase.Calls = nil

			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := httptest.NewRequest("DELETE", tt.path, nil)
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

	controller := controllers.NewCrimeController(
		mockCreateUseCase,
		mockListUseCase,
		mockUpdateStatusUseCase,
		mockDeleteUseCase,
		mockGetStatsUseCase,
	)

	router.GET("/crimes/stats", controller.GetCrimeStats)

	tests := []struct {
		name           string
		mockSetup      func()
		expectedStatus int
		expectedStats  *entities.CrimeStats
	}{
		{
			name: "obtener estadísticas exitosamente",
			mockSetup: func() {
				mockGetStatsUseCase.On("Execute", mock.Anything, mock.Anything).Return(
					&entities.CrimeStats{
						TotalCrimes:      10,
						ActiveCrimes:     8,
						InactiveCrimes:   2,
						CrimesByType:     map[string]int{"ROBO": 5, "ASALTO": 3, "OTRO": 2},
						CrimesByStatus:   map[string]int{"ACTIVE": 8, "INACTIVE": 2},
						CrimesByLocation: map[string]int{"NORTE": 4, "SUR": 3, "ESTE": 2, "OESTE": 1},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedStats: &entities.CrimeStats{
				TotalCrimes:      10,
				ActiveCrimes:     8,
				InactiveCrimes:   2,
				CrimesByType:     map[string]int{"ROBO": 5, "ASALTO": 3, "OTRO": 2},
				CrimesByStatus:   map[string]int{"ACTIVE": 8, "INACTIVE": 2},
				CrimesByLocation: map[string]int{"NORTE": 4, "SUR": 3, "ESTE": 2, "OESTE": 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := httptest.NewRequest("GET", "/crimes/stats", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response entities.CrimeStats
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStats.TotalCrimes, response.TotalCrimes)
				assert.Equal(t, tt.expectedStats.ActiveCrimes, response.ActiveCrimes)
				assert.Equal(t, tt.expectedStats.InactiveCrimes, response.InactiveCrimes)
				assert.Equal(t, tt.expectedStats.CrimesByType, response.CrimesByType)
				assert.Equal(t, tt.expectedStats.CrimesByStatus, response.CrimesByStatus)
				assert.Equal(t, tt.expectedStats.CrimesByLocation, response.CrimesByLocation)
			}
		})
	}
}
