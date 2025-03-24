package main

import (
	"log"

	"go-crime_map_backend/internal/infrastructure/database"
	"go-crime_map_backend/internal/infrastructure/repositories"
	"go-crime_map_backend/internal/interface/controllers"
	"go-crime_map_backend/internal/interface/routes"
	"go-crime_map_backend/internal/usecases"

	"github.com/gin-gonic/gin"
)

func main() {
	// Configurar modo de Gin
	gin.SetMode(gin.DebugMode)

	// Crear el router
	router := gin.Default()

	// Configurar middleware básico
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configurar rutas básicas
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Configurar la base de datos
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer db.Close()

	// Crear el repositorio
	crimeRepo := repositories.NewPostgresCrimeRepository(db)

	// Crear los casos de uso
	createUseCase := usecases.NewCreateCrimeUseCase(crimeRepo)
	listUseCase := usecases.NewListCrimesUseCase(crimeRepo)
	updateStatusUseCase := usecases.NewUpdateCrimeStatusUseCase(crimeRepo)
	deleteUseCase := usecases.NewDeleteCrimeUseCase(crimeRepo)
	getStatsUseCase := usecases.NewGetCrimeStatsUseCase(crimeRepo)
	getCrimeUseCase := usecases.NewGetCrimeUseCase(crimeRepo)

	// Crear el controlador
	crimeController := controllers.NewCrimeController(
		createUseCase,
		listUseCase,
		updateStatusUseCase,
		deleteUseCase,
		getStatsUseCase,
		getCrimeUseCase,
	)

	// Configurar rutas de delitos
	routes.SetupCrimeRoutes(router, crimeController)

	// Imprimir rutas registradas
	log.Println("\nRutas registradas:")
	for _, route := range router.Routes() {
		log.Printf("%s %s\n", route.Method, route.Path)
	}
	log.Println()

	// Iniciar el servidor
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
