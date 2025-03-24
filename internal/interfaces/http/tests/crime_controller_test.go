package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-crime_map_backend/internal/infrastructure/database"
	"go-crime_map_backend/internal/infrastructure/repositories"
	crimeController "go-crime_map_backend/internal/interfaces/http"
	"go-crime_map_backend/internal/usecases"
)

// TestCase representa un caso de prueba para la creación de delitos
type TestCase struct {
	Name             string                 `json:"name"`
	Payload          map[string]interface{} `json:"payload"`
	ExpectedStatus   int                    `json:"expected_status"`
	ShouldValidateDB bool                   `json:"should_validate_db"`
	ExpectedDBState  struct {
		Count       int `json:"count"`
		FirstRecord *struct {
			Type        string `json:"type"`
			Description string `json:"description"`
			Location    struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
				Address   string  `json:"address"`
			} `json:"location"`
		} `json:"first_record,omitempty"`
	} `json:"expected_db_state"`
}

// TestCases representa la estructura del archivo JSON de casos de prueba
type TestCases struct {
	Cases []TestCase `json:"test_cases"`
}

func loadTestCases(t *testing.T) []TestCase {
	data, err := os.ReadFile("testdata/create_crime_test_cases.json")
	require.NoError(t, err, "Error al leer el archivo de casos de prueba")

	var testCases TestCases
	err = json.Unmarshal(data, &testCases)
	require.NoError(t, err, "Error al deserializar los casos de prueba")

	return testCases.Cases
}

func setupTestDB(t *testing.T) *database.Config {
	config := database.NewTestConfig()

	db, err := database.NewPostgresDB(config)
	require.NoError(t, err)
	require.NoError(t, db.Ping())

	// Crear el esquema de pruebas si no existe
	_, err = db.Exec(`CREATE SCHEMA IF NOT EXISTS test`)
	require.NoError(t, err)

	// Eliminar tablas si existen
	_, err = db.Exec(`
		DROP TABLE IF EXISTS test.crimes CASCADE;
		DROP TABLE IF EXISTS test.locations CASCADE;
	`)
	require.NoError(t, err)

	// Leer y ejecutar el esquema SQL
	schemaSQL, err := os.ReadFile("../../../infrastructure/database/schema.sql")
	require.NoError(t, err)

	// Modificar el esquema SQL para usar el esquema de pruebas
	schemaSQL = []byte(fmt.Sprintf("SET search_path TO test;\n%s", string(schemaSQL)))

	_, err = db.Exec(string(schemaSQL))
	require.NoError(t, err)

	return config
}

func processDateString(dateStr string) time.Time {
	switch dateStr {
	case "+24h":
		return time.Now().Add(24 * time.Hour)
	case "-24h":
		return time.Now().Add(-24 * time.Hour)
	default:
		t, _ := time.Parse(time.RFC3339, dateStr)
		return t
	}
}

func TestCreateCrime(t *testing.T) {
	// Cargar casos de prueba
	testCases := loadTestCases(t)

	// Configurar base de datos de prueba
	config := setupTestDB(t)
	db, err := database.NewPostgresDB(config)
	require.NoError(t, err)
	defer db.Close()

	repo := repositories.NewPostgresCrimeRepository(db)

	// Crear controlador
	createCrimeUseCase := usecases.NewCreateCrimeUseCase(repo)
	controller := crimeController.NewCrimeController(createCrimeUseCase)

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			// Limpiar base de datos antes de cada prueba
			err := repo.DeleteAll()
			require.NoError(t, err)

			// Si es el caso de inserción duplicada, crear primero el delito original
			if tt.Name == "error - inserción duplicada" {
				// Procesar la fecha en el payload
				if dateStr, ok := tt.Payload["date"].(string); ok {
					tt.Payload["date"] = processDateString(dateStr).Format(time.RFC3339)
				}

				// Crear el delito original
				payload, err := json.Marshal(tt.Payload)
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/crimes", bytes.NewBuffer(payload))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = req
				controller.Create(c)
				assert.Equal(t, http.StatusCreated, w.Code)
			}

			// Procesar la fecha en el payload
			if dateStr, ok := tt.Payload["date"].(string); ok {
				tt.Payload["date"] = processDateString(dateStr).Format(time.RFC3339)
			}

			// Crear request
			payload, err := json.Marshal(tt.Payload)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/crimes", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			// Crear response recorder
			w := httptest.NewRecorder()

			// Configurar Gin
			gin.SetMode(gin.TestMode)
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Ejecutar handler
			controller.Create(c)

			// Verificar status code
			assert.Equal(t, tt.ExpectedStatus, w.Code)

			// Verificar estado de la base de datos si es necesario
			if tt.ShouldValidateDB {
				crimes, err := repo.GetAll(context.Background())
				require.NoError(t, err)
				assert.Len(t, crimes, tt.ExpectedDBState.Count)

				if tt.ExpectedDBState.FirstRecord != nil && len(crimes) > 0 {
					crime := crimes[0]
					assert.NotEmpty(t, crime.ID)
					assert.Equal(t, tt.ExpectedDBState.FirstRecord.Type, crime.Type)
					assert.Equal(t, tt.ExpectedDBState.FirstRecord.Description, crime.Description)
					assert.Equal(t, tt.ExpectedDBState.FirstRecord.Location.Latitude, crime.Location.Latitude)
					assert.Equal(t, tt.ExpectedDBState.FirstRecord.Location.Longitude, crime.Location.Longitude)
					assert.Equal(t, tt.ExpectedDBState.FirstRecord.Location.Address, crime.Location.Address)
				}
			}
		})
	}
}
