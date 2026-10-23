package usecases

import (
	"context"

	"go-crime_map_backend/internal/domain/entities"
)

// CreateCrimeUseCaseInterface define la interfaz para el caso de uso de creación de delitos
type CreateCrimeUseCaseInterface interface {
	Execute(ctx context.Context, input CreateCrimeInput) (*entities.Crime, error)
}

// ListCrimesUseCaseInterface define la interfaz para el caso de uso de listado de delitos
type ListCrimesUseCaseInterface interface {
	Execute(ctx context.Context, input ListCrimesInput) (*ListCrimesOutput, error)
}

// UpdateCrimeStatusUseCaseInterface define la interfaz para el caso de uso de actualización de estado
type UpdateCrimeStatusUseCaseInterface interface {
	Execute(ctx context.Context, input UpdateCrimeStatusInput) error
}

// DeleteCrimeUseCaseInterface define la interfaz para el caso de uso de eliminación de delitos
type DeleteCrimeUseCaseInterface interface {
	Execute(ctx context.Context, id string) error
}

// GetCrimeStatsUseCaseInterface define la interfaz para el caso de uso de obtención de estadísticas
type GetCrimeStatsUseCaseInterface interface {
	Execute(ctx context.Context, input GetCrimeStatsInput) (*entities.CrimeStats, error)
}

// GetCrimeUseCaseInterface define la interfaz para el caso de uso de obtención de un delito por ID
type GetCrimeUseCaseInterface interface {
	Execute(ctx context.Context, id string) (*entities.Crime, error)
}
