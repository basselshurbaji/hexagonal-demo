package user

import (
	"context"

	"hexagonal-demo/modules/user/internal/entity"
	"hexagonal-demo/modules/user/internal/ports"
)

// ErrNotFound is re-exported from the internal ports package so code outside
// the module can check it with errors.Is.
var ErrNotFound = ports.ErrNotFound

// User is the module's public model. Other modules depend on this, never on
// the internal entity — the mapping keeps the hexagon free to evolve.
type User struct {
	ID    uint64
	Name  string
	Email string
}

// Playlist is the user module's public view of a playlist linked to a user.
type Playlist struct {
	ID   uint64
	Name string
}

// Facade is the module's public interface. Consumers — driving adapters and
// other modules' adapters — depend on this instead of *Module, so they can be
// tested against a fake and the module keeps its internals free to change.
type Facade interface {
	// GetUserByID returns ErrNotFound when no user matches.
	GetUserByID(ctx context.Context, id uint64) (User, error)

	// GetUserPlaylists returns the playlists linked to the user. Returns
	// ErrNotFound when the user does not exist.
	GetUserPlaylists(ctx context.Context, userID uint64) ([]Playlist, error)
}

// facade implements Facade. It holds the module pointer rather than the
// service so it can be handed out before Initialize runs (two-phase setup);
// the service is only dereferenced per call.
type facade struct {
	m *Module
}

var _ Facade = (*facade)(nil)

func fromEntity(e entity.User) User {
	return User{
		ID:    e.ID,
		Name:  e.Name,
		Email: e.Email,
	}
}

func (f *facade) GetUserByID(ctx context.Context, id uint64) (User, error) {
	u, err := f.m.svc.GetUserByID(ctx, id)
	if err != nil {
		return User{}, err
	}
	return fromEntity(u), nil
}

func (f *facade) GetUserPlaylists(ctx context.Context, userID uint64) ([]Playlist, error) {
	playlists, err := f.m.svc.GetUserPlaylists(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]Playlist, 0, len(playlists))
	for _, p := range playlists {
		out = append(out, Playlist{ID: p.ID, Name: p.Name})
	}
	return out, nil
}
