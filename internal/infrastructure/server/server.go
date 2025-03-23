package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	crimeHttp "go-crime_map_backend/internal/interfaces/http"
	"go-crime_map_backend/internal/usecases"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	router     *gin.Engine
}

func NewServer() *Server {
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

	// TODO: Inicializar el repositorio y el caso de uso
	// Por ahora usamos nil como placeholder
	createCrimeUseCase := usecases.NewCreateCrimeUseCase(nil)
	crimeController := crimeHttp.NewCrimeController(createCrimeUseCase)

	// Configurar rutas de delitos
	crimes := router.Group("/api/v1/crimes")
	{
		crimes.POST("/", crimeController.Create)
	}

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
