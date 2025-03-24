package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go-crime_map_backend/internal/infrastructure/database"
	"go-crime_map_backend/internal/infrastructure/repositories"
	"go-crime_map_backend/internal/interface/controllers"
	"go-crime_map_backend/internal/interface/routes"
	"go-crime_map_backend/internal/usecases"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	router     *gin.Engine
}

func NewServer() *Server {
	// Configurar modo de Gin
	gin.SetMode(gin.DebugMode)

	router := gin.Default()

	// Configurar middleware básico
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configurar rutas básicas
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Inicializar la base de datos
	db, err := database.NewPostgresDB()
	if err != nil {
		panic(fmt.Sprintf("Error al conectar con la base de datos: %v", err))
	}

	// Inicializar el repositorio
	crimeRepo := repositories.NewPostgresCrimeRepository(db)

	// Inicializar los casos de uso
	createCrimeUseCase := usecases.NewCreateCrimeUseCase(crimeRepo)
	listCrimesUseCase := usecases.NewListCrimesUseCase(crimeRepo)
	updateStatusUseCase := usecases.NewUpdateCrimeStatusUseCase(crimeRepo)
	deleteCrimeUseCase := usecases.NewDeleteCrimeUseCase(crimeRepo)
	getStatsUseCase := usecases.NewGetCrimeStatsUseCase(crimeRepo)
	getCrimeUseCase := usecases.NewGetCrimeUseCase(crimeRepo)

	crimeController := controllers.NewCrimeController(
		createCrimeUseCase,
		listCrimesUseCase,
		updateStatusUseCase,
		deleteCrimeUseCase,
		getStatsUseCase,
		getCrimeUseCase,
	)

	// Configurar rutas de delitos
	routes.SetupCrimeRoutes(router, crimeController)

	// Imprimir rutas registradas
	fmt.Println("\nRutas registradas:")
	for _, route := range router.Routes() {
		fmt.Printf("%s %s\n", route.Method, route.Path)
	}
	fmt.Println()

	return &Server{
		router: router,
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	fmt.Println("Servidor iniciado en el puerto 8080")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
