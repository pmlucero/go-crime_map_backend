package usecases

import (
	"context"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"

	"github.com/stretchr/testify/assert"
)

func TestListCrimesUseCase_Execute(t *testing.T) {
	// Crear mock del repositorio
	mockRepo := new(MockCrimeRepository)
	useCase := NewListCrimesUseCase(mockRepo)

	// Crear contexto
	ctx := context.Background()

	// Crear datos de prueba
	crimes := []entities.Crime{
		{
			ID:          "1",
			Type:        "ROBO",
			Description: "Robo a mano armada",
			Location: entities.Location{
				Latitude:  -34.603722,
				Longitude: -58.381592,
			},
			Status:    "ACTIVO",
			CreatedAt: time.Now(),
		},
		{
			ID:          "2",
			Type:        "ASALTO",
			Description: "Asalto a comercio",
			Location: entities.Location{
				Latitude:  -34.608722,
				Longitude: -58.382592,
			},
			Status:    "RESUELTO",
			CreatedAt: time.Now(),
		},
	}

	tests := []struct {
		name          string
		input         ListCrimesParams
		mockSetup     func()
		expectedError error
		expectedCount int
	}{
		{
			name: "Listar todos los delitos",
			input: ListCrimesParams{
				Page:  1,
				Limit: 10,
			},
			mockSetup: func() {
				mockRepo.On("List", ctx, 1, 10, (*time.Time)(nil), (*time.Time)(nil), (*string)(nil), (*string)(nil)).Return(crimes, int64(2), nil)
			},
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name: "Listar delitos por tipo",
			input: ListCrimesParams{
				Page:  1,
				Limit: 10,
				Type:  stringPtr("ROBO"),
			},
			mockSetup: func() {
				crimeType := "ROBO"
				mockRepo.On("List", ctx, 1, 10, (*time.Time)(nil), (*time.Time)(nil), &crimeType, (*string)(nil)).Return([]entities.Crime{crimes[0]}, int64(1), nil)
			},
			expectedError: nil,
			expectedCount: 1,
		},
		{
			name: "Listar delitos por rango de fechas",
			input: ListCrimesParams{
				Page:      1,
				Limit:     10,
				StartDate: timePtr(time.Now().Add(-24 * time.Hour)),
				EndDate:   timePtr(time.Now()),
			},
			mockSetup: func() {
				startDate := time.Now().Add(-24 * time.Hour)
				endDate := time.Now()
				mockRepo.On("List", ctx, 1, 10, &startDate, &endDate, (*string)(nil), (*string)(nil)).Return(crimes, int64(2), nil)
			},
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name: "Error en rango de fechas inválido",
			input: ListCrimesParams{
				Page:      1,
				Limit:     10,
				StartDate: timePtr(time.Now()),
				EndDate:   timePtr(time.Now().Add(-24 * time.Hour)),
			},
			mockSetup:     func() {},
			expectedError: ErrInvalidDateRange,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar el mock
			tt.mockSetup()

			// Ejecutar el caso de uso
			result, err := useCase.Execute(ctx, tt.input)

			// Verificar el resultado
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedCount, len(result.Crimes))
			assert.Equal(t, int64(tt.expectedCount), result.Total)

			// Verificar que se llamó al mock
			mockRepo.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
