package user

import (
	"hexagonal-demo/modules/user/internal/ports"
	"hexagonal-demo/modules/user/internal/service"
)

type Module struct {
	svc *service.Service
}

// Initialize assembles the module from its driven ports. Two-phase setup:
// the composition root creates all module references first, then initializes
// them — so adapters can capture module pointers regardless of dependency cycles.
func (m *Module) Initialize(repository ports.Repository, playlists ports.PlaylistCatalog) {
	m.svc = service.New(repository, playlists)
}

// Facade returns the module's public facade. Safe to call before Initialize —
// the facade captures the stable module pointer, not the service.
func (m *Module) Facade() Facade {
	return &facade{m: m}
}
