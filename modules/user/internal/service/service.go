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

// GetUserPlaylists returns the playlists linked to the user. The playlist ids
// are the user module's own data (user_playlists); the playlists themselves
// are hydrated through the PlaylistCatalog port.
func (s *Service) GetUserPlaylists(ctx context.Context, userID uint64) ([]entity.Playlist, error) {
	if _, err := s.repository.GetUserByID(ctx, userID); err != nil {
		return nil, err
	}
	ids, err := s.repository.GetUserPlaylistIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}
	return s.playlists.GetPlaylistsByIDs(ctx, ids)
}
