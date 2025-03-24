package usecases

import "errors"

var (
	// ErrInvalidDateRange indica que el rango de fechas es inválido
	ErrInvalidDateRange = errors.New("el rango de fechas es inválido")

	// ErrInvalidRadius indica que el radio especificado es inválido
	ErrInvalidRadius = errors.New("el radio debe ser mayor a 0")

	// ErrCrimeAlreadyDeleted indica que se intentó modificar un delito que ya fue eliminado
	ErrCrimeAlreadyDeleted = errors.New("el delito ya fue eliminado")

	// ErrInvalidStatusTransition indica que se intentó realizar una transición de estado inválida
	ErrInvalidStatusTransition = errors.New("la transición de estado no es válida")
)
