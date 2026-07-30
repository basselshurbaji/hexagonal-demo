package service

import (
	"context"

	"hexagonal-demo/modules/playlist/internal/entity"
	"hexagonal-demo/modules/playlist/internal/ports"
)

type Service struct {
	repository ports.Repository
	songs      ports.SongCatalog
}

func New(repository ports.Repository, songs ports.SongCatalog) *Service {
	return &Service{repository: repository, songs: songs}
}

func (s *Service) CreatePlaylist(ctx context.Context, name string, userID uint64) (entity.Playlist, error) {
	return s.repository.CreatePlaylist(ctx, name, userID)
}

// AddSongs validates that every song exists (via the SongCatalog port) before
// linking them to the playlist. Returns ports.ErrSongNotFound if any is missing.
func (s *Service) AddSongs(ctx context.Context, playlistID uint64, songIDs []uint64) error {
	if _, err := s.repository.GetPlaylistByID(ctx, playlistID); err != nil {
		return err
	}

	unique := dedupe(songIDs)
	songs, err := s.songs.GetSongsByIDs(ctx, unique)
	if err != nil {
		return err
	}
	if len(songs) != len(unique) {
		return ports.ErrSongNotFound
	}

	return s.repository.AddSongs(ctx, playlistID, unique)
}

func (s *Service) GetPlaylist(ctx context.Context, id uint64) (entity.Playlist, []entity.Song, error) {
	playlist, err := s.repository.GetPlaylistByID(ctx, id)
	if err != nil {
		return entity.Playlist{}, nil, err
	}

	songIDs, err := s.repository.GetSongIDs(ctx, id)
	if err != nil {
		return entity.Playlist{}, nil, err
	}
	if len(songIDs) == 0 {
		return playlist, nil, nil
	}

	songs, err := s.songs.GetSongsByIDs(ctx, songIDs)
	if err != nil {
		return entity.Playlist{}, nil, err
	}
	return playlist, songs, nil
}

func (s *Service) GetPlaylistsByIDs(ctx context.Context, ids []uint64) ([]entity.Playlist, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return s.repository.GetPlaylistsByIDs(ctx, ids)
}

func dedupe(ids []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(ids))
	out := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
