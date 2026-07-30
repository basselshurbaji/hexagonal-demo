package song

import (
	"hexagonal-demo/modules/song/internal/ports"
	"hexagonal-demo/modules/song/internal/service"
)

type Module struct {
	svc *service.Service
}

// Initialize assembles the module from its driven ports. Two-phase setup:
// the composition root creates all module references first, then initializes
// them — so adapters can capture module pointers regardless of dependency cycles.
func (m *Module) Initialize(repository ports.Repository) {
	m.svc = service.New(repository)
}
