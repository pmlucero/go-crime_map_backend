package usecases

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go-crime_map_backend/internal/domain/entities"
	"go-crime_map_backend/internal/domain/repositories"
)

// CrimeStats representa las estadísticas de delitos
type CrimeStats struct {
	TotalCrimes     int                          // Total de delitos
	CrimesByType    map[string]int               // Delitos por tipo
	CrimesByStatus  map[entities.CrimeStatus]int // Delitos por estado
	CrimesByMonth   map[string]int               // Delitos por mes
	RecentCrimes    []*entities.Crime            // Delitos más recientes
	MostCommonTypes []string                     // Tipos de delitos más comunes
}

// GetCrimeStatsInput representa los filtros para obtener estadísticas
type GetCrimeStatsInput struct {
	StartDate time.Time // Fecha inicial
	EndDate   time.Time // Fecha final
	Limit     int       // Límite de resultados para delitos recientes
}

// GetCrimeStatsUseCase implementa la lógica de negocio para obtener estadísticas de delitos
type GetCrimeStatsUseCase struct {
	crimeRepository repositories.CrimeRepository
}

// NewGetCrimeStatsUseCase crea una nueva instancia del caso de uso
func NewGetCrimeStatsUseCase(repo repositories.CrimeRepository) *GetCrimeStatsUseCase {
	return &GetCrimeStatsUseCase{
		crimeRepository: repo,
	}
}

// Execute ejecuta el caso de uso
func (uc *GetCrimeStatsUseCase) Execute(ctx context.Context) (*entities.CrimeStats, error) {
	// Obtener estadísticas del repositorio
	stats, err := uc.crimeRepository.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estadísticas: %w", err)
	}

	return stats, nil
}

// typeCount representa un par de tipo de delito y su cantidad
type typeCount struct {
	crimeType string
	count     int
}

// getMostCommonTypes obtiene los tipos de delitos más comunes
func getMostCommonTypes(crimesByType map[string]int) []string {
	// Crear slice de pares tipo-cantidad
	pairs := make([]typeCount, 0, len(crimesByType))
	for crimeType, count := range crimesByType {
		pairs = append(pairs, typeCount{crimeType, count})
	}

	// Ordenar por cantidad (descendente)
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	// Obtener los 5 tipos más comunes
	result := make([]string, 0, 5)
	for i := 0; i < len(pairs) && i < 5; i++ {
		result = append(result, pairs[i].crimeType)
	}

	return result
}
