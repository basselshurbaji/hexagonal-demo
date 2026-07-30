package ports

import (
	"context"
	"errors"

	"hexagonal-demo/modules/user/internal/entity"
)

// ErrNotFound is returned by Repository implementations when no user matches.
var ErrNotFound = errors.New("user not found")

type Repository interface {
	// GetUserByID returns ErrNotFound when no user exists with the given id.
	GetUserByID(ctx context.Context, id uint64) (entity.User, error)
	// GetUserPlaylistIDs returns the ids of the playlists linked to the user.
	GetUserPlaylistIDs(ctx context.Context, userID uint64) ([]uint64, error)
}
