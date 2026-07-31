package song

import (
	"context"

	"hexagonal-demo/modules/song/internal/entity"
	"hexagonal-demo/modules/song/internal/ports"
)

// ErrNotFound is re-exported from the internal ports package so code outside
// the module can check it with errors.Is.
var ErrNotFound = ports.ErrNotFound

// Song is the module's public model. Other modules depend on this, never on
// the internal entity.
type Song struct {
	ID              uint64
	Name            string
	ArtistID        uint64
	DurationSeconds uint32
	StorageLocation string
}

func fromEntity(e entity.Song) Song {
	return Song{
		ID:              e.ID,
		Name:            e.Name,
		ArtistID:        e.ArtistID,
		DurationSeconds: e.DurationSeconds,
		StorageLocation: e.StorageLocation,
	}
}

// Facade is the module's public interface. Consumers — driving adapters and
// other modules' adapters — depend on this instead of *Module, so they can be
// tested against a fake and the module keeps its internals free to change.
type Facade interface {
	// GetSongsByIDs returns the songs that exist among ids; missing ids are
	// simply omitted from the result.
	GetSongsByIDs(ctx context.Context, ids []uint64) ([]Song, error)
}

// facade implements Facade. It holds the module pointer rather than the
// service so it can be handed out before Initialize runs (two-phase setup);
// the service is only dereferenced per call.
type facade struct {
	m *Module
}

var _ Facade = (*facade)(nil)

func (f *facade) GetSongsByIDs(ctx context.Context, ids []uint64) ([]Song, error) {
	songs, err := f.m.svc.GetSongsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	out := make([]Song, 0, len(songs))
	for _, s := range songs {
		out = append(out, fromEntity(s))
	}
	return out, nil
}
