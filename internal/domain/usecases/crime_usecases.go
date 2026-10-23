package usecases

import (
	"context"
	"time"

	"go-crime_map_backend/internal/domain/entities"
)

// CreateCrimeInput representa los datos de entrada para crear un delito
type CreateCrimeInput struct {
	Title         string
	Description   string
	Type          string
	Latitude      float64
	Longitude     float64
	Address       string
	AddressNumber string
	City          string
	Province      string
	Country       string
	ZipCode       string
}

// CreateCrimeUseCase define la interfaz para crear delitos
type CreateCrimeUseCase interface {
	Execute(ctx context.Context, input CreateCrimeInput) (*entities.Crime, error)
}

// ListCrimesParams representa los parámetros para listar delitos
type ListCrimesParams struct {
	Page      int
	Limit     int
	StartDate *time.Time
	EndDate   *time.Time
	Type      *string
	Status    *string
}

// ListCrimesUseCase define la interfaz para listar delitos
type ListCrimesUseCase interface {
	Execute(ctx context.Context, params ListCrimesParams) (*entities.CrimeList, error)
}

// UpdateCrimeStatusInput representa los datos de entrada para actualizar el estado de un delito
type UpdateCrimeStatusInput struct {
	UUID   string
	Status string
}

// UpdateCrimeStatusUseCase define la interfaz para actualizar el estado de un delito
type UpdateCrimeStatusUseCase interface {
	Execute(ctx context.Context, input UpdateCrimeStatusInput) error
}

// DeleteCrimeUseCase define la interfaz para eliminar delitos
type DeleteCrimeUseCase interface {
	Execute(ctx context.Context, uuid string) error
}

// GetCrimeStatsParams representa los parámetros para obtener estadísticas de delitos
type GetCrimeStatsParams struct {
	StartDate *time.Time
	EndDate   *time.Time
	Type      *string
	Status    *string
}

// GetCrimeStatsUseCase define la interfaz para obtener estadísticas de delitos
type GetCrimeStatsUseCase interface {
	Execute(ctx context.Context) (*entities.CrimeStats, error)
}

// GetCrimeUseCase define la interfaz para obtener un delito por UUID
type GetCrimeUseCase interface {
	Execute(ctx context.Context, uuid string) (*entities.Crime, error)
}
