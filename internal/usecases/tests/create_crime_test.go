package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories/mocks"
	domain_usecases "go-crime_map_backend/internal/domain/usecases"
	"go-crime_map_backend/internal/usecases"
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
			AddressNumber: "1234",
			City:          "Buenos Aires",
			Province:      "CABA",
			Country:       "Argentina",
			ZipCode:       "C1043",
		}
	}

	t.Run("Creación exitosa", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)
		input := createValidInput()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.Crime")).
			Run(func(args mock.Arguments) {
				crime := args.Get(1).(*entities.Crime)
				crime.ID = 1
			}).
			Return(nil)
		// Act
		crime, err := useCase.Execute(ctx, input)
		// Assert
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
	})

	t.Run("Error - título vacío", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)
		input := createValidInput()
		input.Title = ""
		// Act
		crime, err := useCase.Execute(ctx, input)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "el título es requerido", err.Error())
	})

	t.Run("Error - descripción vacía", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)
		input := createValidInput()
		input.Description = ""
		// Act
		crime, err := useCase.Execute(ctx, input)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "la descripción es requerida", err.Error())
	})

	t.Run("Error - tipo vacío", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)
		input := createValidInput()
		input.Type = ""
		// Act
		crime, err := useCase.Execute(ctx, input)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "el tipo es requerido", err.Error())
	})

	t.Run("Error - latitud inválida", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)
		input := createValidInput()
		input.Latitude = 91
		// Act
		crime, err := useCase.Execute(ctx, input)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "la latitud debe estar entre -90 y 90", err.Error())
	})

	t.Run("Error - longitud inválida", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)
		input := createValidInput()
		input.Longitude = 181
		// Act
		crime, err := useCase.Execute(ctx, input)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "la longitud debe estar entre -180 y 180", err.Error())
	})

	t.Run("Error - dirección vacía", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)
		input := createValidInput()
		input.Address = ""
		// Act
		crime, err := useCase.Execute(ctx, input)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "la dirección es requerida", err.Error())
	})

	t.Run("Error - número de dirección nulo", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)
		input := createValidInput()
		input.AddressNumber = ""
		// Act
		crime, err := useCase.Execute(ctx, input)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "el número de la dirección es requerido", err.Error())
	})

	t.Run("Error - ciudad nula", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)
		input := createValidInput()
		input.City = ""
		// Act
		crime, err := useCase.Execute(ctx, input)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "la ciudad es requerida", err.Error())
	})

	t.Run("Error - provincia nula", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)
		input := createValidInput()
		input.Province = ""
		// Act
		crime, err := useCase.Execute(ctx, input)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "la provincia es requerida", err.Error())
	})

	t.Run("Error - país nulo", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)
		input := createValidInput()
		input.Country = ""
		// Act
		crime, err := useCase.Execute(ctx, input)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "el país es requerido", err.Error())
	})

	t.Run("Error - código postal nulo", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)
		input := createValidInput()
		input.ZipCode = ""
		// Act
		crime, err := useCase.Execute(ctx, input)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Equal(t, "el código postal es requerido", err.Error())
	})

	t.Run("Error - error del repositorio", func(t *testing.T) {
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)

		input := createValidInput()

		expectedError := errors.New("error al crear el delito")
		mockRepo.EXPECT().Create(ctx, mock.AnythingOfType("*entities.Crime")).Return(expectedError)

		crime, err := useCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, crime)
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("Error al crear el delito", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)
		input := createValidInput()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.Crime")).Return(assert.AnError)
		// Act
		crime, err := useCase.Execute(ctx, input)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, crime)
	})

	t.Run("Validación inválida", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCrimeRepository(t)
		useCase := usecases.NewCreateCrimeUseCase(mockRepo)
		input := createValidInput()
		input.Title = "" // Inválido
		// Act
		crime, err := useCase.Execute(ctx, input)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, crime)
	})
}
