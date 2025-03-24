package repositories

import "errors"

var (
	// ErrNotFound indica que no se encontró el recurso solicitado
	ErrNotFound = errors.New("recurso no encontrado")
)

// ErrCrimeNotFound indica que el delito no fue encontrado
var ErrCrimeNotFound = errors.New("el delito no existe")
