package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go-crime_map_backend/internal/infrastructure/repositories"
	"go-crime_map_backend/internal/interface/controllers"
	"go-crime_map_backend/internal/interface/routes"
	"go-crime_map_backend/internal/usecases"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Server representa el servidor HTTP
type Server struct {
	httpServer *http.Server
	router     *gin.Engine
}

// NewServer crea una nueva instancia del servidor
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

	return &Server{
		router: router,
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: router,
		},
	}
}

// Start inicia el servidor
func (s *Server) Start() error {
	fmt.Println("Servidor iniciado en el puerto 8080")
	return s.httpServer.ListenAndServe()
}

// Shutdown detiene el servidor
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

// SetupRoutes configura las rutas del servidor
func (s *Server) SetupRoutes(db *sqlx.DB) error {
	// Crear repositorio
	crimeRepo := repositories.NewPostgresCrimeRepository(db)

	// Crear casos de uso
	createCrimeUseCase := usecases.NewCreateCrimeUseCase(crimeRepo)
	listCrimesUseCase := usecases.NewListCrimesUseCase(crimeRepo)
	updateCrimeStatusUseCase := usecases.NewUpdateCrimeStatusUseCase(crimeRepo)
	deleteCrimeUseCase := usecases.NewDeleteCrimeUseCase(crimeRepo)
	getCrimeStatsUseCase := usecases.NewGetCrimeStatsUseCase(crimeRepo)
	getCrimeUseCase := usecases.NewGetCrimeUseCase(crimeRepo)

	// Crear controlador
	crimeController := controllers.NewCrimeController(
		createCrimeUseCase,
		listCrimesUseCase,
		updateCrimeStatusUseCase,
		deleteCrimeUseCase,
		getCrimeStatsUseCase,
		getCrimeUseCase,
	)

	// Configurar rutas
	routes.SetupCrimeRoutes(s.router, crimeController)

	// Imprimir rutas registradas
	fmt.Println("\nRutas registradas:")
	for _, route := range s.router.Routes() {
		fmt.Printf("%s %s\n", route.Method, route.Path)
	}
	fmt.Println()

	return nil
}
