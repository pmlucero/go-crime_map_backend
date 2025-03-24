package http

import (
	"errors"
	"net/http"
	"time"

	"go-crime_map_backend/internal/usecases"

	"github.com/gin-gonic/gin"
)

// CrimeController maneja las peticiones HTTP relacionadas con los delitos
type CrimeController struct {
	createCrimeUseCase *usecases.CreateCrimeUseCase
}

// NewCrimeController crea una nueva instancia del controlador
func NewCrimeController(createCrimeUseCase *usecases.CreateCrimeUseCase) *CrimeController {
	return &CrimeController{
		createCrimeUseCase: createCrimeUseCase,
	}
}

// CreateCrimeRequest representa la estructura de la petición HTTP
type CreateCrimeRequest struct {
	Type        string    `json:"type" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    Location  `json:"location" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
}

// Location representa la ubicación en la petición HTTP
type Location struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Address   string  `json:"address" binding:"required"`
}

// Create maneja la petición POST para crear un nuevo delito
func (c *CrimeController) Create(ctx *gin.Context) {
	var req CreateCrimeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := usecases.CreateCrimeInput{
		Type:        req.Type,
		Description: req.Description,
		Location: usecases.Location{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
			Address:   req.Location.Address,
		},
		Date: req.Date,
	}

	crime, err := c.createCrimeUseCase.Execute(ctx.Request.Context(), input)
	if err != nil {
		var statusCode int
		switch {
		case errors.Is(err, usecases.ErrInvalidType):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecases.ErrEmptyDescription):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecases.ErrDescriptionTooLong):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecases.ErrFutureDate):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecases.ErrInvalidLatitude):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecases.ErrInvalidLongitude):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecases.ErrDuplicateCrime):
			statusCode = http.StatusConflict
		default:
			statusCode = http.StatusInternalServerError
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, crime)
}
