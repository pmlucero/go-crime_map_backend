package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories/mocks"
	"go-crime_map_backend/internal/usecases"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCrimeUseCase_Execute(t *testing.T) {
	// Arrange
	activeCrime := &entities.Crime{
		ID:     1,
		UUID:   uuid.New().String(),
		Status: "ACTIVO",
	}

	deletedCrime := &entities.Crime{
		ID:        2,
		UUID:      uuid.New().String(),
		Status:    "DELETED",
		DeletedAt: &time.Time{},
	}

	type test struct {
		name     string
		uuid     string
		wantErr  bool
		errMsg   string
		mockFunc func(repo *mocks.CrimeRepository)
	}

	tests := []test{
		{
			name:    "Error - UUID vacío",
			uuid:    "",
			wantErr: true,
			errMsg:  "el UUID es requerido",
			mockFunc: func(repo *mocks.CrimeRepository) {
				// No se espera ninguna llamada al repositorio
			},
		},
		{
			name:    "Error - error al obtener el delito",
			uuid:    activeCrime.UUID,
			wantErr: true,
			errMsg:  "error al obtener el delito",
			mockFunc: func(repo *mocks.CrimeRepository) {
				repo.EXPECT().GetByUUID(context.Background(), activeCrime.UUID).Return(nil, errors.New("error al obtener el delito"))
			},
		},
		{
			name:    "Error - delito no encontrado",
			uuid:    activeCrime.UUID,
			wantErr: true,
			errMsg:  "delito no encontrado",
			mockFunc: func(repo *mocks.CrimeRepository) {
				repo.EXPECT().GetByUUID(context.Background(), activeCrime.UUID).Return(nil, nil)
			},
		},
		{
			name:    "Error - delito ya eliminado",
			uuid:    deletedCrime.UUID,
			wantErr: true,
			errMsg:  "el delito ya fue eliminado",
			mockFunc: func(repo *mocks.CrimeRepository) {
				repo.EXPECT().GetByUUID(context.Background(), deletedCrime.UUID).Return(deletedCrime, nil)
			},
		},
		{
			name:    "Error - error al eliminar el delito",
			uuid:    activeCrime.UUID,
			wantErr: true,
			errMsg:  "error al eliminar el delito",
			mockFunc: func(repo *mocks.CrimeRepository) {
				repo.EXPECT().GetByUUID(context.Background(), activeCrime.UUID).Return(activeCrime, nil)
				repo.EXPECT().Delete(context.Background(), activeCrime.ID).Return(errors.New("error al eliminar el delito"))
			},
		},
		{
			name:    "Eliminar delito exitosamente",
			uuid:    activeCrime.UUID,
			wantErr: false,
			errMsg:  "",
			mockFunc: func(repo *mocks.CrimeRepository) {
				repo.EXPECT().GetByUUID(context.Background(), activeCrime.UUID).Return(activeCrime, nil)
				repo.EXPECT().Delete(context.Background(), activeCrime.ID).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := mocks.NewCrimeRepository(t)
			tt.mockFunc(mockRepo)

			useCase := usecases.NewDeleteCrimeUseCase(mockRepo)

			// Act
			err := useCase.Execute(context.Background(), tt.uuid)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
