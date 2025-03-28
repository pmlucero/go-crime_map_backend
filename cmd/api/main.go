package main

import (
	"log"
	"os"

	"go-crime_map_backend/internal/infrastructure/server"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

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
