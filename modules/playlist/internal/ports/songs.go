package ports

import (
	"context"

	"hexagonal-demo/modules/playlist/internal/entity"
)

// SongCatalog is the playlist module's driven port for song data. The playlist
// core depends on this interface only — how songs are actually fetched (song
// module facade today, a remote service tomorrow) is an adapter concern.
type SongCatalog interface {
	// GetSongsByIDs returns the songs that exist among ids; missing ids are omitted.
	GetSongsByIDs(ctx context.Context, ids []uint64) ([]entity.Song, error)
}
