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
	Title         string  `json:"title" binding:"required" example:"Robo a mano armada"`
	Description   string  `json:"description" binding:"required" example:"Robo a mano armada en comercio"`
	Type          string  `json:"type" binding:"required" example:"ROBO"`
	Latitude      float64 `json:"latitude" binding:"required" example:"-34.603722"`
	Longitude     float64 `json:"longitude" binding:"required" example:"-58.381592"`
	Address       string  `json:"address" binding:"required" example:"Av. Corrientes"`
	AddressNumber string  `json:"address_number" binding:"required" example:"1234"`
	City          string  `json:"city" binding:"required" example:"Buenos Aires"`
	Province      string  `json:"province" binding:"required" example:"Buenos Aires"`
	Country       string  `json:"country" binding:"required" example:"Argentina"`
	ZipCode       string  `json:"zip_code" binding:"required" example:"1000"`
}

// @Summary      Crear un nuevo delito
// @Description  Crea un nuevo registro de delito en el sistema
// @Tags         crimes
// @Accept       json
// @Produce      json
// @Param        crime body CreateCrimeRequest true "Datos del delito"
// @Success      201  {object}  entities.Crime
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /crimes [post]
// @Security     ApiKeyAuth
func (c *CrimeController) CreateCrime(ctx *gin.Context) {
	var req CreateCrimeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	input := usecases.CreateCrimeInput{
		Title:         req.Title,
		Description:   req.Description,
		Type:          req.Type,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		Address:       req.Address,
		AddressNumber: req.AddressNumber,
		City:          req.City,
		Province:      req.Province,
		Country:       req.Country,
		ZipCode:       req.ZipCode,
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
	Type      string `form:"type" example:"ROBO"`
	Status    string `form:"status" example:"ACTIVE"`
	StartDate string `form:"start_date" example:"2024-01-01"`
	EndDate   string `form:"end_date" example:"2024-12-31"`
	Limit     int    `form:"limit,default=10" example:"10"`
	Page      int    `form:"page,default=1" example:"1"`
}

// ListCrimesResponse representa la respuesta paginada de delitos
type ListCrimesResponse struct {
	Data       []entities.Crime `json:"data"`
	TotalCount int64            `json:"total_count"`
	Page       int              `json:"page"`
	TotalPages int              `json:"total_pages"`
}

// @Summary      Listar delitos
// @Description  Obtiene una lista paginada de delitos con filtros opcionales
// @Tags         crimes
// @Accept       json
// @Produce      json
// @Param        type query string false "Tipo de delito"
// @Param        status query string false "Estado del delito"
// @Param        start_date query string false "Fecha de inicio"
// @Param        end_date query string false "Fecha de fin"
// @Param        limit query int false "Límite de resultados por página" default(10)
// @Param        page query int false "Número de página" default(1)
// @Success      200  {object}  ListCrimesResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /crimes [get]
// @Security     ApiKeyAuth
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
	Status string `json:"status" binding:"required" example:"INACTIVE"`
}

// @Summary      Actualizar estado de un delito
// @Description  Actualiza el estado de un delito existente
// @Tags         crimes
// @Accept       json
// @Produce      json
// @Param        id path string true "ID del delito"
// @Param        status body UpdateStatusRequest true "Nuevo estado"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /crimes/{id}/status [patch]
// @Security     ApiKeyAuth
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

// @Summary      Eliminar un delito
// @Description  Elimina un delito del sistema
// @Tags         crimes
// @Accept       json
// @Produce      json
// @Param        id path string true "ID del delito"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /crimes/{id} [delete]
// @Security     ApiKeyAuth
func (c *CrimeController) DeleteCrime(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.deleteUseCase.Execute(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary      Obtener estadísticas de delitos
// @Description  Obtiene estadísticas generales sobre los delitos en el sistema
// @Tags         crimes
// @Accept       json
// @Produce      json
// @Success      200  {object}  entities.CrimeStats
// @Failure      500  {object}  ErrorResponse
// @Router       /crimes/stats [get]
// @Security     ApiKeyAuth
func (c *CrimeController) GetCrimeStats(ctx *gin.Context) {
	stats, err := c.getStatsUseCase.Execute(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// @Summary      Obtener un delito por ID
// @Description  Obtiene los detalles de un delito específico
// @Tags         crimes
// @Accept       json
// @Produce      json
// @Param        id path string true "ID del delito"
// @Success      200  {object}  entities.Crime
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /crimes/{id} [get]
// @Security     ApiKeyAuth
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
	Error string `json:"error" example:"Error message"`
}

// Helper functions
func parseFloat64(s string) float64 {
	if s == "" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
