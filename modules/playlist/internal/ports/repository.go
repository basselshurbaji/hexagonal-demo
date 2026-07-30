package ports

import (
	"context"
	"errors"

	"hexagonal-demo/modules/playlist/internal/entity"
)

// ErrNotFound is returned by Repository implementations when no playlist matches.
var ErrNotFound = errors.New("playlist not found")

// ErrSongNotFound is returned by the service when adding songs that do not exist.
var ErrSongNotFound = errors.New("one or more songs do not exist")

type Repository interface {
	// CreatePlaylist creates a playlist and links it to the given user.
	CreatePlaylist(ctx context.Context, name string, userID uint64) (entity.Playlist, error)
	// GetPlaylistByID returns ErrNotFound when no playlist exists with the given id.
	GetPlaylistByID(ctx context.Context, id uint64) (entity.Playlist, error)
	// AddSongs links the given songs to the playlist; already-linked songs are ignored.
	AddSongs(ctx context.Context, playlistID uint64, songIDs []uint64) error
	// GetSongIDs returns the ids of the songs in the playlist.
	GetSongIDs(ctx context.Context, playlistID uint64) ([]uint64, error)
	// GetPlaylistsByIDs returns the playlists that exist among ids; missing ids
	// are simply omitted from the result — no error.
	GetPlaylistsByIDs(ctx context.Context, ids []uint64) ([]entity.Playlist, error)
}
