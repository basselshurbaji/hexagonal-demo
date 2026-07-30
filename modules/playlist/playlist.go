package playlist

import (
	"hexagonal-demo/modules/playlist/internal/ports"
	"hexagonal-demo/modules/playlist/internal/service"
)

type Module struct {
	svc *service.Service
}

// Initialize assembles the module from its driven ports. Two-phase setup:
// the composition root creates all module references first, then initializes
// them — so adapters can capture module pointers regardless of dependency cycles.
func (m *Module) Initialize(repository ports.Repository, songs ports.SongCatalog) {
	m.svc = service.New(repository, songs)
}
