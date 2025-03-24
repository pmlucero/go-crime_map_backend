package repositories

import "errors"

var (
	// ErrNotFound indica que no se encontró el recurso solicitado
	ErrNotFound = errors.New("recurso no encontrado")
)
