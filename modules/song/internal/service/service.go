package service

import (
	"context"

	"hexagonal-demo/modules/song/internal/entity"
	"hexagonal-demo/modules/song/internal/ports"
)

type Service struct {
	repository ports.Repository
}

func New(repository ports.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetSongsByIDs(ctx context.Context, ids []uint64) ([]entity.Song, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return s.repository.GetSongsByIDs(ctx, ids)
}
