package repositories

import (
	"context"
	"sync"

	"go-crime_map_backend/internal/domain/entities"
)

// MemoryCrimeRepository implementa el repositorio de delitos en memoria
type MemoryCrimeRepository struct {
	mu     sync.RWMutex
	crimes map[string]*entities.Crime
}

// NewMemoryCrimeRepository crea una nueva instancia del repositorio en memoria
func NewMemoryCrimeRepository() *MemoryCrimeRepository {
	return &MemoryCrimeRepository{
		crimes: make(map[string]*entities.Crime),
	}
}

// Create guarda un nuevo delito en el repositorio
func (r *MemoryCrimeRepository) Create(ctx context.Context, crime *entities.Crime) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.crimes[crime.ID] = crime
	return nil
}

// GetByID obtiene un delito por su ID
func (r *MemoryCrimeRepository) GetByID(ctx context.Context, id string) (*entities.Crime, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if crime, exists := r.crimes[id]; exists {
		return crime, nil
	}
	return nil, nil
}

// GetAll obtiene todos los delitos
func (r *MemoryCrimeRepository) GetAll(ctx context.Context) ([]*entities.Crime, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	crimes := make([]*entities.Crime, 0, len(r.crimes))
	for _, crime := range r.crimes {
		crimes = append(crimes, crime)
	}
	return crimes, nil
}

// Update actualiza un delito existente
func (r *MemoryCrimeRepository) Update(ctx context.Context, crime *entities.Crime) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.crimes[crime.ID]; exists {
		r.crimes[crime.ID] = crime
		return nil
	}
	return nil
}

// Delete elimina un delito por su ID
func (r *MemoryCrimeRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.crimes, id)
	return nil
}
