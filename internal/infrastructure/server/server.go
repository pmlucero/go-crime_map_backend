package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"go-crime_map_backend/internal/infrastructure/database"
	"go-crime_map_backend/internal/infrastructure/repositories"
	crimeHttp "go-crime_map_backend/internal/interfaces/http"
	"go-crime_map_backend/internal/usecases"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	router     *gin.Engine
	db         *sql.DB
}

// NewServer crea una nueva instancia del servidor HTTP
func NewServer() *Server {
	router := gin.Default()

	// Inicializar la conexión a la base de datos
	var dbConfig *database.Config
	if os.Getenv("TEST_MODE") == "true" {
		fmt.Println("Usando configuración de base de datos de test")
		dbConfig = database.NewTestConfig()
	} else {
		fmt.Println("Usando configuración de base de datos de producción")
		dbConfig = database.NewConfig()
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		panic(fmt.Sprintf("Error al conectar con la base de datos: %v", err))
	}

	// Inicializar el repositorio de PostgreSQL
	crimeRepo := repositories.NewPostgresCrimeRepository(db)

	// Inicializar el caso de uso
	createCrimeUseCase := usecases.NewCreateCrimeUseCase(crimeRepo)

	// Inicializar el controlador
	crimeController := crimeHttp.NewCrimeController(createCrimeUseCase)

	// Configurar rutas
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	// Grupo de rutas para la API v1
	v1 := router.Group("/api/v1")
	{
		crimes := v1.Group("/crimes")
		{
			crimes.POST("/", crimeController.Create)
		}
	}

	return &Server{
		router: router,
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: router,
		},
		db: db,
	}
}

func (s *Server) Start() error {
	fmt.Println("Servidor iniciado en el puerto 8080")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Cerrar la conexión a la base de datos
	if err := s.db.Close(); err != nil {
		fmt.Printf("Error al cerrar la conexión a la base de datos: %v\n", err)
	}

	return s.httpServer.Shutdown(ctx)
}
