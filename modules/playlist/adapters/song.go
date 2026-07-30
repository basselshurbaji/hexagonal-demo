package adapters

import (
	"context"

	"hexagonal-demo/modules/playlist/internal/entity"
	"hexagonal-demo/modules/playlist/internal/ports"
	"hexagonal-demo/modules/song"
)

// compile-time check: SongModuleAdapter satisfies the SongCatalog port
var _ ports.SongCatalog = (*SongModuleAdapter)(nil)

// SongModuleAdapter implements the SongCatalog port by calling the song
// module's facade, translating song.Song into the playlist module's own model.
type SongModuleAdapter struct {
	songs *song.Module
}

func NewSongModuleAdapter(songs *song.Module) *SongModuleAdapter {
	return &SongModuleAdapter{songs: songs}
}

func (a *SongModuleAdapter) GetSongsByIDs(ctx context.Context, ids []uint64) ([]entity.Song, error) {
	res, err := a.songs.GetSongsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	songs := make([]entity.Song, 0, len(res))
	for _, s := range res {
		songs = append(songs, entity.Song{
			ID:              s.ID,
			Name:            s.Name,
			ArtistID:        s.ArtistID,
			DurationSeconds: s.DurationSeconds,
			StorageLocation: s.StorageLocation,
		})
	}
	return songs, nil
}
