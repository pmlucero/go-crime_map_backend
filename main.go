package main

import (
	"database/sql"
	"log"

	"go-crime_map_backend/internal/infrastructure/repositories"
	infrausecases "go-crime_map_backend/internal/infrastructure/usecases"
	"go-crime_map_backend/internal/interface/controllers"
	"go-crime_map_backend/internal/interface/middleware"
	"go-crime_map_backend/internal/security"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Inicializar la conexión a la base de datos
	db, err := sql.Open("postgres", "postgres://user:password@localhost/crime_map?sslmode=disable")
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error cerrando la base de datos: %v", err)
		}
	}()

	// Inicializar los repositorios
	crimeRepo := repositories.NewCrimeRepository(db)
	securityRepo := security.NewRepository(db)

	// Crear los casos de uso
	createUseCase := infrausecases.NewCreateCrimeUseCase(crimeRepo)
	listUseCase := infrausecases.NewListCrimesUseCase(crimeRepo)
	updateStatusUseCase := infrausecases.NewUpdateCrimeStatusUseCase(crimeRepo)
	deleteUseCase := infrausecases.NewDeleteCrimeUseCase(crimeRepo)
	getStatsUseCase := infrausecases.NewGetCrimeStatsUseCase(crimeRepo)
	getCrimeUseCase := infrausecases.NewGetCrimeUseCase(crimeRepo)

	// Crear el router
	router := gin.Default()

	// Configurar el middleware de autenticación
	router.Use(middleware.AuthMiddleware(securityRepo))

	// Crear los controladores
	crimeController := controllers.NewCrimeController(
		createUseCase,
		listUseCase,
		updateStatusUseCase,
		deleteUseCase,
		getStatsUseCase,
		getCrimeUseCase,
	)

	// Configurar las rutas
	router.POST("/crimes", crimeController.CreateCrime)
	router.GET("/crimes", crimeController.ListCrimes)
	router.PATCH("/crimes/:id/status", crimeController.UpdateCrimeStatus)
	router.DELETE("/crimes/:id", crimeController.DeleteCrime)
	router.GET("/crimes/stats", crimeController.GetCrimeStats)
	router.GET("/crimes/:id", crimeController.GetCrime)

	// Iniciar el servidor
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
