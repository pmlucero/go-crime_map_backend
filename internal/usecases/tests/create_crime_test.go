package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	domain_usecases "go-crime_map_backend/internal/domain/usecases"
	"go-crime_map_backend/internal/mocks"
	"go-crime_map_backend/internal/usecases"
	"go-crime_map_backend/internal/utils"
)

func TestCreateCrimeUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	// Función helper para crear un input válido base
	createValidInput := func() domain_usecases.CreateCrimeInput {
		return domain_usecases.CreateCrimeInput{
			Title:         "Robo",
			Description:   "Robo en tienda",
			Type:          "ROBBERY",
			Latitude:      -34.603722,
			Longitude:     -58.381592,
			Address:       "Av. Corrientes",
			AddressNumber: utils.StringPtr("1234"),
			City:          utils.StringPtr("Buenos Aires"),
			Province:      utils.StringPtr("CABA"),
			Country:       utils.StringPtr("Argentina"),
			ZipCode:       utils.StringPtr("C1043"),
		}
	}

	t.Run("Creación exitosa", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)

		input := createValidInput()

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.Crime")).Return(nil)

		crime, err := useCase.Execute(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, crime)
		assert.NotEmpty(t, crime.ID)
		assert.Equal(t, input.Title, crime.Title)
		assert.Equal(t, input.Description, crime.Description)
		assert.Equal(t, input.Type, crime.Type)
		assert.Equal(t, "ACTIVE", crime.Status)
		assert.Equal(t, input.Latitude, crime.Location.Latitude)
		assert.Equal(t, input.Longitude, crime.Location.Longitude)
		assert.Equal(t, input.Address, crime.Location.Address)
		assert.Equal(t, input.AddressNumber, crime.Location.AddressNumber)
		assert.Equal(t, input.City, crime.Location.City)
		assert.Equal(t, input.Province, crime.Location.Province)
		assert.Equal(t, input.Country, crime.Location.Country)
		assert.Equal(t, input.ZipCode, crime.Location.ZipCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - título vacío", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)

		input := createValidInput()
		input.Title = ""

		crime, err := useCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "el título es requerido", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error - descripción vacía", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)

		input := createValidInput()
		input.Description = ""

		crime, err := useCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "la descripción es requerida", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error - tipo vacío", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)

		input := createValidInput()
		input.Type = ""

		crime, err := useCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "el tipo es requerido", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error - latitud inválida", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)

		input := createValidInput()
		input.Latitude = 91

		crime, err := useCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "la latitud debe estar entre -90 y 90", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error - longitud inválida", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)

		input := createValidInput()
		input.Longitude = 181

		crime, err := useCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "la longitud debe estar entre -180 y 180", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error - dirección vacía", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)

		input := createValidInput()
		input.Address = ""

		crime, err := useCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "la dirección es requerida", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error - número de dirección nulo", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)

		input := createValidInput()
		input.AddressNumber = nil

		crime, err := useCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "el número de la dirección es requerido", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error - ciudad nula", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)

		input := createValidInput()
		input.City = nil

		crime, err := useCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "la ciudad es requerida", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error - provincia nula", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)

		input := createValidInput()
		input.Province = nil

		crime, err := useCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "la provincia es requerida", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error - país nulo", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)

		input := createValidInput()
		input.Country = nil

		crime, err := useCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "el país es requerido", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error - código postal nulo", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)

		input := createValidInput()
		input.ZipCode = nil

		crime, err := useCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "el código postal es requerido", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error - error del repositorio", func(t *testing.T) {
		mockRepo := new(mocks.MockCrimeRepository)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)

		input := createValidInput()

		expectedError := errors.New("error al crear el delito")
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.Crime")).Return(expectedError)

		crime, err := useCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Contains(t, err.Error(), expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}
