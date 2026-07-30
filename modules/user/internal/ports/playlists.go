package ports

import (
	"context"

	"hexagonal-demo/modules/user/internal/entity"
)

// PlaylistCatalog is the user module's driven port for playlist data. The
// user↔playlist association is owned by the playlist side of the port — this
// module only asks questions in terms of users.
type PlaylistCatalog interface {
	// GetPlaylistsByUserID returns the playlists linked to the given user.
	GetPlaylistsByUserID(ctx context.Context, userID uint64) ([]entity.Playlist, error)
}
