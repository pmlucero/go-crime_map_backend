package tests

import (
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"

	"github.com/stretchr/testify/assert"
)

// coverage: 100%
func TestCrime_Validate(t *testing.T) {
	// Arrange
	validUUID := "123e4567-e89b-4456-8842-456612345678"
	now := time.Now()

	tests := []struct {
		name    string
		crime   entities.Crime
		wantErr bool
		errMsg  string
	}{
		{
			name: "Crime válido",
			crime: entities.Crime{
				UUID:        validUUID,
				Type:        string(entities.CrimeTypeRobo),
				Description: "Descripción del delito",
				Location: entities.Location{
					Latitude:  -34.603722,
					Longitude: -58.381592,
				},
				Status:    string(entities.CrimeStatusActive),
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: false,
		},
		{
			name: "UUID vacío",
			crime: entities.Crime{
				Type:        string(entities.CrimeTypeRobo),
				Description: "Descripción del delito",
				Location: entities.Location{
					Latitude:  -34.603722,
					Longitude: -58.381592,
				},
				Status:    string(entities.CrimeStatusActive),
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: true,
			errMsg:  "el UUID es requerido",
		},
		{
			name: "UUID inválido",
			crime: entities.Crime{
				UUID:        "uuid-invalido",
				Type:        string(entities.CrimeTypeRobo),
				Description: "Descripción del delito",
				Location: entities.Location{
					Latitude:  -34.603722,
					Longitude: -58.381592,
				},
				Status:    string(entities.CrimeStatusActive),
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: true,
			errMsg:  "UUID inválido",
		},
		{
			name: "Latitud inválida",
			crime: entities.Crime{
				UUID:        validUUID,
				Type:        string(entities.CrimeTypeRobo),
				Description: "Descripción del delito",
				Location: entities.Location{
					Latitude:  91.0, // Inválida
					Longitude: -58.381592,
				},
				Status:    string(entities.CrimeStatusActive),
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: true,
			errMsg:  "la latitud debe estar entre -90 y 90",
		},
		{
			name: "Longitud inválida",
			crime: entities.Crime{
				UUID:        validUUID,
				Type:        string(entities.CrimeTypeRobo),
				Description: "Descripción del delito",
				Location: entities.Location{
					Latitude:  -34.603722,
					Longitude: 181.0, // Inválida
				},
				Status:    string(entities.CrimeStatusActive),
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: true,
			errMsg:  "la longitud debe estar entre -180 y 180",
		},
		{
			name: "Tipo de delito inválido",
			crime: entities.Crime{
				UUID:        validUUID,
				Type:        "TIPO_INVALIDO",
				Description: "Descripción del delito",
				Location: entities.Location{
					Latitude:  -34.603722,
					Longitude: -58.381592,
				},
				Status:    string(entities.CrimeStatusActive),
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: true,
			errMsg:  "tipo de delito inválido",
		},
		{
			name: "Descripción vacía",
			crime: entities.Crime{
				UUID: validUUID,
				Type: string(entities.CrimeTypeRobo),
				Location: entities.Location{
					Latitude:  -34.603722,
					Longitude: -58.381592,
				},
				Status:    string(entities.CrimeStatusActive),
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: true,
			errMsg:  "la descripción es requerida",
		},
		{
			name: "Estado inválido",
			crime: entities.Crime{
				UUID:        validUUID,
				Type:        string(entities.CrimeTypeRobo),
				Description: "Descripción del delito",
				Location: entities.Location{
					Latitude:  -34.603722,
					Longitude: -58.381592,
				},
				Status:    "ESTADO_INVALIDO",
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: true,
			errMsg:  "estado inválido",
		},
		{
			name: "Fecha de actualización anterior a la fecha de creación",
			crime: entities.Crime{
				UUID:        validUUID,
				Type:        string(entities.CrimeTypeRobo),
				Description: "Descripción del delito",
				Location: entities.Location{
					Latitude:  -34.603722,
					Longitude: -58.381592,
				},
				Status:    string(entities.CrimeStatusActive),
				CreatedAt: now,
				UpdatedAt: now.Add(-time.Hour), // Una hora antes
			},
			wantErr: true,
			errMsg:  "la fecha de actualización no puede ser anterior a la fecha de creación",
		},
	}

	// Act
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			crime := tt.crime
			// Act
			err := crime.Validate()
			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
