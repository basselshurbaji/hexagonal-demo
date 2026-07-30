package ports

import (
	"context"
	"errors"

	"hexagonal-demo/modules/song/internal/entity"
)

// ErrNotFound is returned by Repository implementations when no song matches.
var ErrNotFound = errors.New("song not found")

type Repository interface {
	// GetSongsByIDs returns the songs that exist among ids; missing ids are
	// simply omitted from the result — no error.
	GetSongsByIDs(ctx context.Context, ids []uint64) ([]entity.Song, error)
}
