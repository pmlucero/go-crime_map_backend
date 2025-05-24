package routes

import (
	"go-crime_map_backend/internal/interface/controllers"

	"github.com/gin-gonic/gin"
)

func SetupCrimeRoutes(router *gin.Engine, controller *controllers.CrimeController, authMiddleware gin.HandlerFunc) {
	protected := router.Group("/")
	protected.Use(authMiddleware)
	{
		protected.POST("/crimes", controller.CreateCrime)
		protected.GET("/crimes", controller.ListCrimes)
		protected.GET("/crimes/stats", controller.GetCrimeStats)
		protected.GET("/crimes/:id", controller.GetCrime)
		protected.PATCH("/crimes/:id/status", controller.UpdateCrimeStatus)
		protected.DELETE("/crimes", controller.DeleteCrime)
		protected.DELETE("/crimes/:id", controller.DeleteCrime)
	}
}
