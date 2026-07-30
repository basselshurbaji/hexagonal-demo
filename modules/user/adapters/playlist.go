package adapters

import (
	"context"

	"hexagonal-demo/modules/playlist"
	"hexagonal-demo/modules/user/internal/entity"
	"hexagonal-demo/modules/user/internal/ports"
)

// compile-time check: PlaylistModuleAdapter satisfies the PlaylistCatalog port
var _ ports.PlaylistCatalog = (*PlaylistModuleAdapter)(nil)

// PlaylistModuleAdapter implements the PlaylistCatalog port by calling the
// playlist module's facade, translating playlist.Playlist into the user
// module's own model.
type PlaylistModuleAdapter struct {
	playlists *playlist.Module
}

func NewPlaylistModuleAdapter(playlists *playlist.Module) *PlaylistModuleAdapter {
	return &PlaylistModuleAdapter{playlists: playlists}
}

func (a *PlaylistModuleAdapter) GetPlaylistsByIDs(ctx context.Context, ids []uint64) ([]entity.Playlist, error) {
	res, err := a.playlists.GetPlaylistsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	playlists := make([]entity.Playlist, 0, len(res))
	for _, p := range res {
		playlists = append(playlists, entity.Playlist{ID: p.ID, Name: p.Name})
	}
	return playlists, nil
}
