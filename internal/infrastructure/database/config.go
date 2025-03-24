package database

import (
	"fmt"
	"os"
)

// Config representa la configuración de la base de datos
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Schema   string
}

// NewConfig crea una nueva configuración de base de datos
func NewConfig() *Config {
	return &Config{
		Host:     "",                // Host vacío para usar socket Unix
		Port:     "",                // Puerto vacío para usar socket Unix
		User:     os.Getenv("USER"), // Usuario actual del sistema
		Password: "",
		DBName:   "crime_map",
		SSLMode:  "disable",
		Schema:   "public",
	}
}

// NewTestConfig crea una nueva configuración de base de datos para pruebas
func NewTestConfig() *Config {
	return &Config{
		Host:     "",                // Host vacío para usar socket Unix
		Port:     "",                // Puerto vacío para usar socket Unix
		User:     os.Getenv("USER"), // Usuario actual del sistema
		Password: "",
		DBName:   "crime_map_test",
		SSLMode:  "disable",
		Schema:   "test",
	}
}

// DSN retorna la cadena de conexión para PostgreSQL
func (c *Config) DSN() string {
	if c.Host == "" && c.Port == "" {
		// Usar socket Unix
		return fmt.Sprintf("dbname=%s sslmode=%s search_path=%s",
			c.DBName, c.SSLMode, c.Schema)
	}
	// Fallback a TCP si se especifican host y puerto
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode, c.Schema)
}

// getEnvOrDefault obtiene una variable de entorno o retorna un valor por defecto
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
