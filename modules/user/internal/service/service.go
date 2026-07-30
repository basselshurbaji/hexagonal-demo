package service

import (
	"context"

	"hexagonal-demo/modules/user/internal/entity"
	"hexagonal-demo/modules/user/internal/ports"
)

type Service struct {
	repository ports.Repository
	playlists  ports.PlaylistCatalog
}

func New(repository ports.Repository, playlists ports.PlaylistCatalog) *Service {
	return &Service{repository: repository, playlists: playlists}
}

func (s *Service) GetUserByID(ctx context.Context, id uint64) (entity.User, error) {
	return s.repository.GetUserByID(ctx, id)
}

// GetUserPlaylists returns the playlists linked to the user, resolved
// entirely through the PlaylistCatalog port — the user module has no
// knowledge of how the association is stored.
func (s *Service) GetUserPlaylists(ctx context.Context, userID uint64) ([]entity.Playlist, error) {
	if _, err := s.repository.GetUserByID(ctx, userID); err != nil {
		return nil, err
	}
	return s.playlists.GetPlaylistsByUserID(ctx, userID)
}
