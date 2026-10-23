package usecases

import (
	"context"
	"errors"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
	"go-crime_map_backend/internal/domain/usecases"
)

// CreateCrimeUseCaseImpl implementa el caso de uso de creación de delitos
type CreateCrimeUseCaseImpl struct {
	repo repositories.CrimeRepository
}

// NewCreateCrimeUseCase crea una nueva instancia del caso de uso
func NewCreateCrimeUseCase(repo repositories.CrimeRepository) usecases.CreateCrimeUseCase {
	return &CreateCrimeUseCaseImpl{repo: repo}
}

func (uc *CreateCrimeUseCaseImpl) Execute(ctx context.Context, input usecases.CreateCrimeInput) (*entities.Crime, error) {
	crime := &entities.Crime{
		Title:       input.Title,
		Description: input.Description,
		Type:        input.Type,
		Status:      "ACTIVE",
		Location: entities.Location{
			Latitude:      input.Latitude,
			Longitude:     input.Longitude,
			Address:       input.Address,
			AddressNumber: input.AddressNumber,
			City:          input.City,
			Province:      input.Province,
			Country:       input.Country,
			ZipCode:       input.ZipCode,
		},
	}

	if err := uc.repo.Create(ctx, crime); err != nil {
		return nil, err
	}

	return crime, nil
}

// ListCrimesUseCaseImpl implementa el caso de uso de listado de delitos
type ListCrimesUseCaseImpl struct {
	repo repositories.CrimeRepository
}

// NewListCrimesUseCase crea una nueva instancia del caso de uso
func NewListCrimesUseCase(repo repositories.CrimeRepository) usecases.ListCrimesUseCase {
	return &ListCrimesUseCaseImpl{repo: repo}
}

func (uc *ListCrimesUseCaseImpl) Execute(ctx context.Context, params usecases.ListCrimesParams) (*entities.CrimeList, error) {
	crimes, total, err := uc.repo.List(ctx, params.Page, params.Limit, params.StartDate, params.EndDate, params.Type, params.Status)
	if err != nil {
		return nil, err
	}

	return &entities.CrimeList{
		Items: crimes,
		Total: total,
	}, nil
}

// UpdateCrimeStatusUseCaseImpl implementa el caso de uso de actualización de estado
type UpdateCrimeStatusUseCaseImpl struct {
	repo repositories.CrimeRepository
}

// NewUpdateCrimeStatusUseCase crea una nueva instancia del caso de uso
func NewUpdateCrimeStatusUseCase(repo repositories.CrimeRepository) usecases.UpdateCrimeStatusUseCase {
	return &UpdateCrimeStatusUseCaseImpl{repo: repo}
}

func (uc *UpdateCrimeStatusUseCaseImpl) Execute(ctx context.Context, input usecases.UpdateCrimeStatusInput) error {
	crime, err := uc.repo.GetByUUID(ctx, input.UUID)
	if err != nil {
		return err
	}

	if crime == nil {
		return errors.New("crime not found")
	}

	crime.Status = input.Status
	return uc.repo.Update(ctx, crime)
}

// DeleteCrimeUseCaseImpl implementa el caso de uso de eliminación de delitos
type DeleteCrimeUseCaseImpl struct {
	repo repositories.CrimeRepository
}

// NewDeleteCrimeUseCase crea una nueva instancia del caso de uso
func NewDeleteCrimeUseCase(repo repositories.CrimeRepository) usecases.DeleteCrimeUseCase {
	return &DeleteCrimeUseCaseImpl{repo: repo}
}

func (uc *DeleteCrimeUseCaseImpl) Execute(ctx context.Context, uuid string) error {
	crime, err := uc.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	if crime == nil {
		return errors.New("crime not found")
	}

	return uc.repo.Delete(ctx, crime.ID)
}

// GetCrimeStatsUseCaseImpl implementa el caso de uso de obtención de estadísticas
type GetCrimeStatsUseCaseImpl struct {
	repo repositories.CrimeRepository
}

// NewGetCrimeStatsUseCase crea una nueva instancia del caso de uso
func NewGetCrimeStatsUseCase(repo repositories.CrimeRepository) usecases.GetCrimeStatsUseCase {
	return &GetCrimeStatsUseCaseImpl{repo: repo}
}

func (uc *GetCrimeStatsUseCaseImpl) Execute(ctx context.Context) (*entities.CrimeStats, error) {
	return uc.repo.GetStats(ctx)
}

// GetCrimeUseCaseImpl implementa el caso de uso de obtención de un delito
type GetCrimeUseCaseImpl struct {
	repo repositories.CrimeRepository
}

// NewGetCrimeUseCase crea una nueva instancia del caso de uso
func NewGetCrimeUseCase(repo repositories.CrimeRepository) usecases.GetCrimeUseCase {
	return &GetCrimeUseCaseImpl{repo: repo}
}

func (uc *GetCrimeUseCaseImpl) Execute(ctx context.Context, uuid string) (*entities.Crime, error) {
	crime, err := uc.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if crime == nil {
		return nil, errors.New("crime not found")
	}

	return crime, nil
}
