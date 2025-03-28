package main

import (
	"log"
	"os"

	"go-crime_map_backend/internal/infrastructure/server"

	_ "go-crime_map_backend/docs" // Importar la documentación de Swagger

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// @title           Crime Map API
// @version         1.0
// @description     API para el sistema de gestión de crímenes Crime Map
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Obtener variables de entorno
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/crime_map?sslmode=disable"
	}

	// Conectar a la base de datos
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// Crear servidor
	srv := server.NewServer()

	// Configurar rutas
	if err := srv.SetupRoutes(db); err != nil {
		log.Fatalf("Error al configurar rutas: %v", err)
	}

	// Iniciar servidor
	if err := srv.Start(); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
