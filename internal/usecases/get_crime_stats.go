package usecases

import (
	"context"
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

// GetCrimeStatsUseCase maneja la lógica de negocio para obtener estadísticas
type GetCrimeStatsUseCase struct {
	crimeRepo repositories.CrimeRepository
}

// NewGetCrimeStatsUseCase crea una nueva instancia del caso de uso
func NewGetCrimeStatsUseCase(repo repositories.CrimeRepository) *GetCrimeStatsUseCase {
	return &GetCrimeStatsUseCase{
		crimeRepo: repo,
	}
}

// Execute ejecuta el caso de uso para obtener estadísticas
func (uc *GetCrimeStatsUseCase) Execute(ctx context.Context, input GetCrimeStatsInput) (*CrimeStats, error) {
	// Validar que las fechas sean coherentes
	if !input.StartDate.IsZero() && !input.EndDate.IsZero() && input.StartDate.After(input.EndDate) {
		return nil, ErrInvalidDateRange
	}

	// Validar que el límite sea positivo
	if input.Limit <= 0 {
		input.Limit = 5 // Valor por defecto
	}

	// Obtener delitos con los filtros especificados
	crimes, err := uc.crimeRepo.List(ctx, repositories.ListCrimesFilter{
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		Limit:     input.Limit,
		Offset:    0,
	})
	if err != nil {
		return nil, err
	}

	// Calcular estadísticas
	stats := &CrimeStats{
		TotalCrimes:    len(crimes),
		CrimesByType:   make(map[string]int),
		CrimesByStatus: make(map[entities.CrimeStatus]int),
		CrimesByMonth:  make(map[string]int),
		RecentCrimes:   crimes,
	}

	// Procesar estadísticas
	for _, crime := range crimes {
		// Contar por tipo
		stats.CrimesByType[crime.Type]++

		// Contar por estado
		stats.CrimesByStatus[crime.Status]++

		// Contar por mes
		monthKey := crime.CreatedAt.Format("2006-01")
		stats.CrimesByMonth[monthKey]++
	}

	// Obtener tipos más comunes
	stats.MostCommonTypes = getMostCommonTypes(stats.CrimesByType)

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
