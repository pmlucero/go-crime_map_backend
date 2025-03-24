package controllers

import (
	"net/http"
	"strconv"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/usecases"

	"github.com/gin-gonic/gin"
)

type CrimeController struct {
	createUseCase       *usecases.CreateCrimeUseCase
	listUseCase         *usecases.ListCrimesUseCase
	updateStatusUseCase *usecases.UpdateCrimeStatusUseCase
	deleteUseCase       *usecases.DeleteCrimeUseCase
	getStatsUseCase     *usecases.GetCrimeStatsUseCase
	getUseCase          *usecases.GetCrimeUseCase
}

func NewCrimeController(
	createUseCase *usecases.CreateCrimeUseCase,
	listUseCase *usecases.ListCrimesUseCase,
	updateStatusUseCase *usecases.UpdateCrimeStatusUseCase,
	deleteUseCase *usecases.DeleteCrimeUseCase,
	getStatsUseCase *usecases.GetCrimeStatsUseCase,
	getUseCase *usecases.GetCrimeUseCase,
) *CrimeController {
	return &CrimeController{
		createUseCase:       createUseCase,
		listUseCase:         listUseCase,
		updateStatusUseCase: updateStatusUseCase,
		deleteUseCase:       deleteUseCase,
		getStatsUseCase:     getStatsUseCase,
		getUseCase:          getUseCase,
	}
}

// CreateCrime godoc
// @Summary Crear un nuevo delito
// @Description Crea un nuevo delito en el sistema
// @Tags crimes
// @Accept json
// @Produce json
// @Param crime body CreateCrimeRequest true "Datos del delito"
// @Success 201 {object} entities.Crime
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /crimes [post]
func (c *CrimeController) CreateCrime(ctx *gin.Context) {
	var req CreateCrimeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := usecases.CreateCrimeInput{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Location: entities.Location{
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
			Address:   req.Address,
		},
	}

	crime, err := c.createUseCase.Execute(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

// ListCrimes godoc
// @Summary Listar delitos
// @Description Obtiene una lista paginada de delitos con filtros opcionales
// @Tags crimes
// @Accept json
// @Produce json
// @Param type query string false "Tipo de delito"
// @Param status query string false "Estado del delito"
// @Param start_date query string false "Fecha inicial (formato: 2006-01-02T15:04:05Z)"
// @Param end_date query string false "Fecha final (formato: 2006-01-02T15:04:05Z)"
// @Param limit query int false "Límite de resultados por página (default: 10, max: 100)"
// @Param page query int false "Número de página (default: 1)"
// @Success 200 {object} ListCrimesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /crimes [get]
func (c *CrimeController) ListCrimes(ctx *gin.Context) {
	// Obtener y validar parámetros de paginación
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Obtener y validar fechas
	var startDate, endDate *time.Time
	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inicial inválido. Use YYYY-MM-DD"})
			return
		}
		startDate = &parsedDate
	}

	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha final inválido. Use YYYY-MM-DD"})
			return
		}
		endDate = &parsedDate
	}

	// Obtener y validar tipo y estado
	var crimeType, status *string
	if typeStr := ctx.Query("type"); typeStr != "" {
		crimeType = &typeStr
	}
	if statusStr := ctx.Query("status"); statusStr != "" {
		status = &statusStr
	}

	// Ejecutar caso de uso
	result, err := c.listUseCase.Execute(ctx.Request.Context(), usecases.ListCrimesParams{
		Page:      page,
		Limit:     limit,
		StartDate: startDate,
		EndDate:   endDate,
		Type:      crimeType,
		Status:    status,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// UpdateCrimeStatus godoc
// @Summary Actualizar estado de un delito
// @Description Actualiza el estado de un delito existente
// @Tags crimes
// @Accept json
// @Produce json
// @Param id path string true "ID del delito"
// @Param status body UpdateStatusRequest true "Nuevo estado"
// @Success 200 {object} entities.Crime
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /crimes/{id}/status [patch]
func (c *CrimeController) UpdateCrimeStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID no proporcionado"})
		return
	}

	var req UpdateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := usecases.UpdateCrimeStatusInput{
		ID:        id,
		NewStatus: req.Status,
	}

	err := c.updateStatusUseCase.Execute(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// DeleteCrime godoc
// @Summary Eliminar un delito
// @Description Realiza una eliminación lógica de un delito
// @Tags crimes
// @Accept json
// @Produce json
// @Param id path string true "ID del delito"
// @Success 200 {object} entities.Crime
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /crimes/{id} [delete]
func (c *CrimeController) DeleteCrime(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID no proporcionado"})
		return
	}

	err := c.deleteUseCase.Execute(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// GetCrimeStats godoc
// @Summary Obtener estadísticas de delitos
// @Description Obtiene estadísticas sobre los delitos registrados
// @Tags crimes
// @Accept json
// @Produce json
// @Success 200 {object} entities.CrimeStats
// @Failure 500 {object} ErrorResponse
// @Router /crimes/stats [get]
func (c *CrimeController) GetCrimeStats(ctx *gin.Context) {
	stats, err := c.getStatsUseCase.Execute(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// GetCrime godoc
// @Summary Obtener un delito por ID
// @Description Obtiene los detalles de un delito específico
// @Tags crimes
// @Accept json
// @Produce json
// @Param id path string true "ID del delito"
// @Success 200 {object} entities.Crime
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /crimes/{id} [get]
func (c *CrimeController) GetCrime(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID no proporcionado"})
		return
	}

	crime, err := c.getUseCase.Execute(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, crime)
}

// Request/Response structs
type CreateCrimeRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Address     string  `json:"address" binding:"required"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

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
