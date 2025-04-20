package tests

import (
	"os"
	"testing"

	"go-crime_map_backend/internal/infrastructure/database"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgresDB(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "Conexión exitosa con valores por defecto",
			envVars: map[string]string{
				"DB_HOST":     "",
				"DB_PORT":     "",
				"DB_USER":     "",
				"DB_PASSWORD": "",
				"DB_NAME":     "",
			},
			wantErr: false,
		},
		{
			name: "Conexión exitosa con valores personalizados",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USER":     "postgres",
				"DB_PASSWORD": "postgres",
				"DB_NAME":     "crime_map",
			},
			wantErr: false,
		},
		{
			name: "Error de conexión - Puerto inválido",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "9999",
				"DB_USER":     "postgres",
				"DB_PASSWORD": "postgres",
				"DB_NAME":     "crime_map",
			},
			wantErr: true,
		},
		{
			name: "Error de conexión - Host inválido",
			envVars: map[string]string{
				"DB_HOST":     "invalid_host",
				"DB_PORT":     "5432",
				"DB_USER":     "postgres",
				"DB_PASSWORD": "postgres",
				"DB_NAME":     "crime_map",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Guardar variables de entorno originales
			originalEnv := make(map[string]string)
			for key := range tt.envVars {
				originalEnv[key] = os.Getenv(key)
			}

			// Establecer variables de entorno para el test
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Restaurar variables de entorno originales al finalizar
			defer func() {
				for key, value := range originalEnv {
					os.Setenv(key, value)
				}
			}()

			// Ejecutar el test
			db, err := database.NewPostgresDB()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				if assert.NoError(t, err) {
					assert.NotNil(t, db)
					// Cerrar la conexión
					err = db.Close()
					assert.NoError(t, err)
				}
			}
		})
	}
}
