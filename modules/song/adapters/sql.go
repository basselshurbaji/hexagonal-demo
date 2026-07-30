package adapters

import (
	"context"

	sqlcgen "hexagonal-demo/db/gen"
	"hexagonal-demo/modules/song/internal/entity"
	"hexagonal-demo/modules/song/internal/ports"
)

// compile-time check: SqlAdapter satisfies the Repository port
var _ ports.Repository = (*SqlAdapter)(nil)

type SqlAdapter struct {
	queries *sqlcgen.Queries
}

func NewSqlAdapter(db sqlcgen.DBTX) *SqlAdapter {
	return &SqlAdapter{queries: sqlcgen.New(db)}
}

func (a *SqlAdapter) GetSongsByIDs(ctx context.Context, ids []uint64) ([]entity.Song, error) {
	rows, err := a.queries.GetSongsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	songs := make([]entity.Song, 0, len(rows))
	for _, row := range rows {
		songs = append(songs, entity.Song{
			ID:              row.ID,
			Name:            row.Name,
			ArtistID:        row.ArtistID,
			DurationSeconds: row.DurationSeconds,
			StorageLocation: row.StorageLocation,
		})
	}
	return songs, nil
}
