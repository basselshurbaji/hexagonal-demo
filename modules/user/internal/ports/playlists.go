package ports

import (
	"context"

	"hexagonal-demo/modules/user/internal/entity"
)

// PlaylistCatalog is the user module's driven port for playlist data. The user
// core resolves which playlist ids belong to a user (its own data) and uses
// this port to hydrate them.
type PlaylistCatalog interface {
	// GetPlaylistsByIDs returns the playlists that exist among ids; missing ids are omitted.
	GetPlaylistsByIDs(ctx context.Context, ids []uint64) ([]entity.Playlist, error)
}
