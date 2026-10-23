package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func NewPostgresDB() (*sql.DB, error) {
	// Obtener variables de entorno
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "crime_map"
	}

	// Construir la cadena de conexión
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// Abrir conexión
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la base de datos: %v", err)
	}

	// Verificar la conexión
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error al conectar con la base de datos: %v", err)
	}

	// Habilitar las extensiones necesarias
	if _, err := db.Exec("CREATE EXTENSION IF NOT EXISTS cube"); err != nil {
		return nil, fmt.Errorf("error al crear la extensión cube: %v", err)
	}

	if _, err := db.Exec("CREATE EXTENSION IF NOT EXISTS earthdistance"); err != nil {
		return nil, fmt.Errorf("error al crear la extensión earthdistance: %v", err)
	}

	return db, nil
}
