package tests

import (
	"testing"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestCrime_Validate(t *testing.T) {
	tests := []struct {
		name    string
		crime   entities.Crime
		wantErr bool
	}{
		{
			name: "Delito válido",
			crime: entities.Crime{
				Title:       "Robo en tienda",
				Description: "Robo en tienda de conveniencia",
				Type:        string(entities.CrimeTypeRobo),
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
					City:          utils.StringPtr("Buenos Aires"),
					Province:      utils.StringPtr("Buenos Aires"),
					Country:       utils.StringPtr("Argentina"),
					ZipCode:       utils.StringPtr("1000"),
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Delito sin título",
			crime: entities.Crime{
				Description: "Robo en tienda de conveniencia",
				Type:        string(entities.CrimeTypeRobo),
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
			},
			wantErr: true,
		},
		{
			name: "Delito sin descripción",
			crime: entities.Crime{
				Title:  "Robo en tienda",
				Type:   string(entities.CrimeTypeRobo),
				Status: string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
			},
			wantErr: true,
		},
		{
			name: "Delito sin tipo",
			crime: entities.Crime{
				Title:       "Robo en tienda",
				Description: "Robo en tienda de conveniencia",
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
			},
			wantErr: true,
		},
		{
			name: "Delito con tipo inválido",
			crime: entities.Crime{
				Title:       "Robo en tienda",
				Description: "Robo en tienda de conveniencia",
				Type:        "TIPO_INVALIDO",
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
			},
			wantErr: true,
		},
		{
			name: "Delito sin estado",
			crime: entities.Crime{
				Title:       "Robo en tienda",
				Description: "Robo en tienda de conveniencia",
				Type:        string(entities.CrimeTypeRobo),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
			},
			wantErr: true,
		},
		{
			name: "Delito con estado inválido",
			crime: entities.Crime{
				Title:       "Robo en tienda",
				Description: "Robo en tienda de conveniencia",
				Type:        string(entities.CrimeTypeRobo),
				Status:      "ESTADO_INVALIDO",
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
			},
			wantErr: true,
		},
		{
			name: "Delito sin ubicación",
			crime: entities.Crime{
				Title:       "Robo en tienda",
				Description: "Robo en tienda de conveniencia",
				Type:        string(entities.CrimeTypeRobo),
				Status:      string(entities.CrimeStatusActive),
			},
			wantErr: true,
		},
		{
			name: "Delito con ubicación inválida",
			crime: entities.Crime{
				Title:       "Robo en tienda",
				Description: "Robo en tienda de conveniencia",
				Type:        string(entities.CrimeTypeRobo),
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Address: "Av. Corrientes",
				},
			},
			wantErr: true,
		},
		{
			name: "Delito con fecha de actualización anterior a creación",
			crime: entities.Crime{
				Title:       "Robo en tienda",
				Description: "Robo en tienda de conveniencia",
				Type:        string(entities.CrimeTypeRobo),
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			wantErr: true,
		},
		{
			name: "Delito con tipo de delito válido (ROBO)",
			crime: entities.Crime{
				Title:       "Robo en tienda",
				Description: "Robo en tienda de conveniencia",
				Type:        string(entities.CrimeTypeRobo),
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
			},
			wantErr: false,
		},
		{
			name: "Delito con tipo de delito válido (AGRESION)",
			crime: entities.Crime{
				Title:       "Agresión en calle",
				Description: "Agresión en la vía pública",
				Type:        string(entities.CrimeTypeAgresion),
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
			},
			wantErr: false,
		},
		{
			name: "Delito con tipo de delito válido (HURTO)",
			crime: entities.Crime{
				Title:       "Hurto en transporte",
				Description: "Hurto en transporte público",
				Type:        string(entities.CrimeTypeHurto),
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
			},
			wantErr: false,
		},
		{
			name: "Delito con tipo de delito válido (VANDALISMO)",
			crime: entities.Crime{
				Title:       "Vandalismo en propiedad",
				Description: "Vandalismo en propiedad privada",
				Type:        string(entities.CrimeTypeVandalismo),
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
			},
			wantErr: false,
		},
		{
			name: "Delito con tipo de delito válido (OTRO)",
			crime: entities.Crime{
				Title:       "Otro tipo de delito",
				Description: "Descripción de otro delito",
				Type:        string(entities.CrimeTypeOtro),
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
			},
			wantErr: false,
		},
		{
			name: "Delito con estado válido (ACTIVE)",
			crime: entities.Crime{
				Title:       "Robo en tienda",
				Description: "Robo en tienda de conveniencia",
				Type:        string(entities.CrimeTypeRobo),
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
			},
			wantErr: false,
		},
		{
			name: "Delito con estado válido (INACTIVE)",
			crime: entities.Crime{
				Title:       "Robo en tienda",
				Description: "Robo en tienda de conveniencia",
				Type:        string(entities.CrimeTypeRobo),
				Status:      string(entities.CrimeStatusInactive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
			},
			wantErr: false,
		},
		{
			name: "Delito con estado válido (DELETED)",
			crime: entities.Crime{
				Title:       "Robo en tienda",
				Description: "Robo en tienda de conveniencia",
				Type:        string(entities.CrimeTypeRobo),
				Status:      string(entities.CrimeStatusDeleted),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
			},
			wantErr: false,
		},
		{
			name: "Delito con fechas válidas",
			crime: entities.Crime{
				Title:       "Robo en tienda",
				Description: "Robo en tienda de conveniencia",
				Type:        string(entities.CrimeTypeRobo),
				Status:      string(entities.CrimeStatusActive),
				Location: entities.Location{
					Latitude:      40.7128,
					Longitude:     -74.0060,
					Address:       "Av. Corrientes",
					AddressNumber: utils.StringPtr("1234"),
				},
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.crime.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLocation_Validate(t *testing.T) {
	tests := []struct {
		name     string
		location entities.Location
		wantErr  bool
	}{
		{
			name: "Ubicación válida",
			location: entities.Location{
				Latitude:      40.7128,
				Longitude:     -74.0060,
				Address:       "Av. Corrientes",
				AddressNumber: utils.StringPtr("1234"),
				City:          utils.StringPtr("Buenos Aires"),
				Province:      utils.StringPtr("Buenos Aires"),
				Country:       utils.StringPtr("Argentina"),
				ZipCode:       utils.StringPtr("1000"),
			},
			wantErr: false,
		},
		{
			name: "Ubicación sin latitud",
			location: entities.Location{
				Longitude:     -74.0060,
				Address:       "Av. Corrientes",
				AddressNumber: utils.StringPtr("1234"),
			},
			wantErr: true,
		},
		{
			name: "Ubicación sin longitud",
			location: entities.Location{
				Latitude:      40.7128,
				Address:       "Av. Corrientes",
				AddressNumber: utils.StringPtr("1234"),
			},
			wantErr: true,
		},
		{
			name: "Ubicación sin dirección",
			location: entities.Location{
				Latitude:      40.7128,
				Longitude:     -74.0060,
				AddressNumber: utils.StringPtr("1234"),
			},
			wantErr: true,
		},
		{
			name: "Ubicación con latitud inválida",
			location: entities.Location{
				Latitude:      200.0,
				Longitude:     -74.0060,
				Address:       "Av. Corrientes",
				AddressNumber: utils.StringPtr("1234"),
			},
			wantErr: true,
		},
		{
			name: "Ubicación con longitud inválida",
			location: entities.Location{
				Latitude:      40.7128,
				Longitude:     200.0,
				Address:       "Av. Corrientes",
				AddressNumber: utils.StringPtr("1234"),
			},
			wantErr: true,
		},
		{
			name: "Ubicación con coordenadas límite válidas",
			location: entities.Location{
				Latitude:      90.0,
				Longitude:     180.0,
				Address:       "Av. Corrientes",
				AddressNumber: utils.StringPtr("1234"),
			},
			wantErr: false,
		},
		{
			name: "Ubicación con coordenadas límite negativas válidas",
			location: entities.Location{
				Latitude:      -90.0,
				Longitude:     -180.0,
				Address:       "Av. Corrientes",
				AddressNumber: utils.StringPtr("1234"),
			},
			wantErr: false,
		},
		{
			name: "Ubicación con coordenadas límite inválidas",
			location: entities.Location{
				Latitude:      91.0,
				Longitude:     181.0,
				Address:       "Av. Corrientes",
				AddressNumber: utils.StringPtr("1234"),
			},
			wantErr: true,
		},
		{
			name: "Ubicación con coordenadas límite negativas inválidas",
			location: entities.Location{
				Latitude:      -91.0,
				Longitude:     -181.0,
				Address:       "Av. Corrientes",
				AddressNumber: utils.StringPtr("1234"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.location.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
