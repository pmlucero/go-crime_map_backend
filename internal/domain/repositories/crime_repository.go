package repositories

import (
	"context"
	"go-crime_map_backend/internal/domain/entities"
)

// CrimeRepository define las operaciones que se pueden realizar con los delitos
type CrimeRepository interface {
	// Create guarda un nuevo delito en el repositorio
	Create(ctx context.Context, crime *entities.Crime) error

	// GetByID obtiene un delito por su ID
	GetByID(ctx context.Context, id string) (*entities.Crime, error)

	// GetAll obtiene todos los delitos
	GetAll(ctx context.Context) ([]*entities.Crime, error)

	// Update actualiza un delito existente
	Update(ctx context.Context, crime *entities.Crime) error

	// Delete elimina un delito por su ID
	Delete(ctx context.Context, id string) error
}
