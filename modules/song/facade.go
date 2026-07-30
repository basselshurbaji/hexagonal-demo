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

// GetSongsByIDs returns the songs that exist among ids; missing ids are
// simply omitted from the result.
func (m *Module) GetSongsByIDs(ctx context.Context, ids []uint64) ([]Song, error) {
	songs, err := m.svc.GetSongsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	out := make([]Song, 0, len(songs))
	for _, s := range songs {
		out = append(out, fromEntity(s))
	}
	return out, nil
}
