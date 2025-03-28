package controllers

import (
	"net/http"
	"strconv"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/usecases"

	"github.com/gin-gonic/gin"
)

// CrimeController maneja las peticiones relacionadas con delitos
type CrimeController struct {
	createUseCase       usecases.CreateCrimeUseCase
	listUseCase         usecases.ListCrimesUseCase
	updateStatusUseCase usecases.UpdateCrimeStatusUseCase
	deleteUseCase       usecases.DeleteCrimeUseCase
	getStatsUseCase     usecases.GetCrimeStatsUseCase
	getCrimeUseCase     usecases.GetCrimeUseCase
}

// NewCrimeController crea una nueva instancia del controlador
func NewCrimeController(
	createUseCase usecases.CreateCrimeUseCase,
	listUseCase usecases.ListCrimesUseCase,
	updateStatusUseCase usecases.UpdateCrimeStatusUseCase,
	deleteUseCase usecases.DeleteCrimeUseCase,
	getStatsUseCase usecases.GetCrimeStatsUseCase,
	getCrimeUseCase usecases.GetCrimeUseCase,
) *CrimeController {
	return &CrimeController{
		createUseCase:       createUseCase,
		listUseCase:         listUseCase,
		updateStatusUseCase: updateStatusUseCase,
		deleteUseCase:       deleteUseCase,
		getStatsUseCase:     getStatsUseCase,
		getCrimeUseCase:     getCrimeUseCase,
	}
}

// CreateCrimeRequest representa la estructura de la petición para crear un delito
type CreateCrimeRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Address     string  `json:"address" binding:"required"`
}

// CreateCrime maneja la creación de un nuevo delito
func (c *CrimeController) CreateCrime(ctx *gin.Context) {
	var req CreateCrimeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	input := usecases.CreateCrimeInput{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Address:     req.Address,
	}

	crime, err := c.createUseCase.Execute(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, crime)
}

// ListCrimesRequest representa la solicitud para listar delitos
type ListCrimesRequest struct {
	Type      string `form:"type"`
	Status    string `form:"status"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Limit     int    `form:"limit,default=10"`
	Page      int    `form:"page,default=1"`
}

// ListCrimesResponse representa la respuesta paginada de delitos
type ListCrimesResponse struct {
	Data       []entities.Crime `json:"data"`
	TotalCount int64            `json:"total_count"`
	Page       int              `json:"page"`
	TotalPages int              `json:"total_pages"`
}

// ListCrimes maneja el listado de delitos
func (c *CrimeController) ListCrimes(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	params := usecases.ListCrimesParams{
		Page:  page,
		Limit: limit,
	}

	crimes, err := c.listUseCase.Execute(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, crimes)
}

// UpdateStatusRequest representa la estructura de la petición para actualizar el estado de un delito
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateCrimeStatus maneja la actualización del estado de un delito
func (c *CrimeController) UpdateCrimeStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	input := usecases.UpdateCrimeStatusInput{
		ID:     id,
		Status: req.Status,
	}

	if err := c.updateStatusUseCase.Execute(ctx, input); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// DeleteCrime maneja la eliminación de un delito
func (c *CrimeController) DeleteCrime(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.deleteUseCase.Execute(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// GetCrimeStats maneja la obtención de estadísticas de delitos
func (c *CrimeController) GetCrimeStats(ctx *gin.Context) {
	stats, err := c.getStatsUseCase.Execute(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// GetCrime maneja la obtención de un delito por ID
func (c *CrimeController) GetCrime(ctx *gin.Context) {
	id := ctx.Param("id")

	crime, err := c.getCrimeUseCase.Execute(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, crime)
}

// ErrorResponse representa la estructura de una respuesta de error
type ErrorResponse struct {
	Error string `json:"error"`
}

// Helper functions
func parseFloat64(s string) float64 {
	if s == "" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
