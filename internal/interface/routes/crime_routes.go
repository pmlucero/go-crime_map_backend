package routes

import (
	"go-crime_map_backend/internal/interface/controllers"

	"github.com/gin-gonic/gin"
)

func SetupCrimeRoutes(router *gin.Engine, controller *controllers.CrimeController) {
	router.POST("/crimes", controller.CreateCrime)
	router.GET("/crimes", controller.ListCrimes)
	router.GET("/crimes/stats", controller.GetCrimeStats)
	router.GET("/crimes/:id", controller.GetCrime)
	router.PATCH("/crimes/:id/status", controller.UpdateCrimeStatus)
	router.DELETE("/crimes", controller.DeleteCrime)
	router.DELETE("/crimes/:id", controller.DeleteCrime)
}
