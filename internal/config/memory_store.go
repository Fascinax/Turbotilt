package config

import (
	"sync"
)

// MemoryStore est une structure qui stocke les services sélectionnés en mémoire
type MemoryStore struct {
	SelectedServices []ManifestService
	mu               sync.RWMutex
}

// memoryStore est l'instance unique qui stocke la sélection en mémoire
var memoryStore *MemoryStore
var once sync.Once

// GetMemoryStore retourne l'instance unique de MemoryStore (singleton)
func GetMemoryStore() *MemoryStore {
	once.Do(func() {
		memoryStore = &MemoryStore{
			SelectedServices: []ManifestService{},
		}
	})
	return memoryStore
}

// StoreSelectedServices stocke les services sélectionnés en mémoire
func (ms *MemoryStore) StoreSelectedServices(services []ManifestService) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.SelectedServices = services
}

// GetSelectedServices retourne les services sélectionnés stockés en mémoire
func (ms *MemoryStore) GetSelectedServices() []ManifestService {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	// Retourner une copie pour éviter les modifications externes
	servicesCopy := make([]ManifestService, len(ms.SelectedServices))
	copy(servicesCopy, ms.SelectedServices)
	return servicesCopy
}

// HasSelectedServices vérifie si des services ont été sélectionnés
func (ms *MemoryStore) HasSelectedServices() bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return len(ms.SelectedServices) > 0
}

// ClearSelectedServices efface les services sélectionnés de la mémoire
func (ms *MemoryStore) ClearSelectedServices() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.SelectedServices = []ManifestService{}
}

// GetManifestFromMemory crée un objet Manifest à partir des services stockés en mémoire
func GetManifestFromMemory() Manifest {
	store := GetMemoryStore()
	return Manifest{
		Services: store.GetSelectedServices(),
	}
}
