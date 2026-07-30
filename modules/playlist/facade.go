package playlist

import (
	"context"

	"hexagonal-demo/modules/playlist/internal/entity"
	"hexagonal-demo/modules/playlist/internal/ports"
)

// ErrNotFound is re-exported from the internal ports package so code outside
// the module can check it with errors.Is.
var ErrNotFound = ports.ErrNotFound

// ErrSongNotFound is returned when adding songs that do not exist.
var ErrSongNotFound = ports.ErrSongNotFound

// Playlist is the module's public model.
type Playlist struct {
	ID   uint64
	Name string
}

// Song is the playlist module's public view of a song in a playlist.
type Song struct {
	ID              uint64
	Name            string
	ArtistID        uint64
	DurationSeconds uint32
	StorageLocation string
}

// PlaylistWithSongs bundles a playlist with its songs.
type PlaylistWithSongs struct {
	Playlist
	Songs []Song
}

func fromEntity(e entity.Playlist) Playlist {
	return Playlist{ID: e.ID, Name: e.Name}
}

// CreatePlaylist creates a playlist owned by the given user.
func (m *Module) CreatePlaylist(ctx context.Context, name string, userID uint64) (Playlist, error) {
	playlist, err := m.svc.CreatePlaylist(ctx, name, userID)
	if err != nil {
		return Playlist{}, err
	}
	return fromEntity(playlist), nil
}

// AddSongs links songs to a playlist. Returns ErrNotFound if the playlist
// does not exist and ErrSongNotFound if any song id does not exist.
func (m *Module) AddSongs(ctx context.Context, playlistID uint64, songIDs []uint64) error {
	return m.svc.AddSongs(ctx, playlistID, songIDs)
}

// GetPlaylistsByIDs returns the playlists that exist among ids; missing ids
// are simply omitted from the result.
func (m *Module) GetPlaylistsByIDs(ctx context.Context, ids []uint64) ([]Playlist, error) {
	playlists, err := m.svc.GetPlaylistsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	out := make([]Playlist, 0, len(playlists))
	for _, p := range playlists {
		out = append(out, fromEntity(p))
	}
	return out, nil
}

// GetPlaylist returns a playlist with its songs. Returns ErrNotFound if the
// playlist does not exist.
func (m *Module) GetPlaylist(ctx context.Context, id uint64) (PlaylistWithSongs, error) {
	playlist, songs, err := m.svc.GetPlaylist(ctx, id)
	if err != nil {
		return PlaylistWithSongs{}, err
	}
	out := PlaylistWithSongs{
		Playlist: fromEntity(playlist),
		Songs:    make([]Song, 0, len(songs)),
	}
	for _, s := range songs {
		out.Songs = append(out.Songs, Song{
			ID:              s.ID,
			Name:            s.Name,
			ArtistID:        s.ArtistID,
			DurationSeconds: s.DurationSeconds,
			StorageLocation: s.StorageLocation,
		})
	}
	return out, nil
}
