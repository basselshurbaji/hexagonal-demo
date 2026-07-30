package adapters

import (
	"context"
	"database/sql"
	"errors"

	sqlcgen "hexagonal-demo/db/gen"
	"hexagonal-demo/modules/playlist/internal/entity"
	"hexagonal-demo/modules/playlist/internal/ports"
)

// compile-time check: SqlAdapter satisfies the Repository port
var _ ports.Repository = (*SqlAdapter)(nil)

type SqlAdapter struct {
	db      *sql.DB
	queries *sqlcgen.Queries
}

func NewSqlAdapter(db *sql.DB) *SqlAdapter {
	return &SqlAdapter{db: db, queries: sqlcgen.New(db)}
}

func (a *SqlAdapter) CreatePlaylist(ctx context.Context, name string, userID uint64) (entity.Playlist, error) {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return entity.Playlist{}, err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := a.queries.WithTx(tx)
	result, err := qtx.CreatePlaylist(ctx, name)
	if err != nil {
		return entity.Playlist{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return entity.Playlist{}, err
	}
	err = qtx.LinkPlaylistToUser(ctx, sqlcgen.LinkPlaylistToUserParams{
		UserID:     userID,
		PlaylistID: uint64(id),
	})
	if err != nil {
		return entity.Playlist{}, err
	}
	if err := tx.Commit(); err != nil {
		return entity.Playlist{}, err
	}
	return entity.Playlist{ID: uint64(id), Name: name}, nil
}

func (a *SqlAdapter) GetPlaylistByID(ctx context.Context, id uint64) (entity.Playlist, error) {
	row, err := a.queries.GetPlaylistByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Playlist{}, ports.ErrNotFound
	}
	if err != nil {
		return entity.Playlist{}, err
	}
	return entity.Playlist{ID: row.ID, Name: row.Name}, nil
}

func (a *SqlAdapter) AddSongs(ctx context.Context, playlistID uint64, songIDs []uint64) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := a.queries.WithTx(tx)
	for _, songID := range songIDs {
		err := qtx.AddSongToPlaylist(ctx, sqlcgen.AddSongToPlaylistParams{
			PlaylistID: playlistID,
			SongID:     songID,
		})
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (a *SqlAdapter) GetSongIDs(ctx context.Context, playlistID uint64) ([]uint64, error) {
	return a.queries.GetPlaylistSongIDs(ctx, playlistID)
}

func (a *SqlAdapter) GetPlaylistsByIDs(ctx context.Context, ids []uint64) ([]entity.Playlist, error) {
	rows, err := a.queries.GetPlaylistsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	playlists := make([]entity.Playlist, 0, len(rows))
	for _, row := range rows {
		playlists = append(playlists, entity.Playlist{ID: row.ID, Name: row.Name})
	}
	return playlists, nil
}

func (a *SqlAdapter) GetPlaylistsByUserID(ctx context.Context, userID uint64) ([]entity.Playlist, error) {
	rows, err := a.queries.GetPlaylistsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	playlists := make([]entity.Playlist, 0, len(rows))
	for _, row := range rows {
		playlists = append(playlists, entity.Playlist{ID: row.ID, Name: row.Name})
	}
	return playlists, nil
}
