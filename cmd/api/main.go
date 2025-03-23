package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-crime_map_backend/internal/infrastructure/server"
)

func main() {
	// Inicializar el servidor
	srv := server.NewServer()
	
	// Iniciar el servidor en una goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	// Esperar señal de interrupción
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Cerrar el servidor de manera elegante
	if err := srv.Shutdown(); err != nil {
		log.Fatalf("Error al cerrar el servidor: %v", err)
	}
} 