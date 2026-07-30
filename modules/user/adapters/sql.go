package adapters

import (
	"context"
	sqlcgen "hexagonal-demo/db/gen"
	"hexagonal-demo/modules/user/internal/entity"
	"hexagonal-demo/modules/user/internal/ports"
)

// compile-time check: SqlAdapter satisfies the Repository port
var _ ports.Repository = (*SqlAdapter)(nil)

type SqlAdapter struct {
	queries *sqlcgen.Queries
}

func NewSqlAdapter(db sqlcgen.DBTX) *SqlAdapter {
	return &SqlAdapter{queries: sqlcgen.New(db)}
}

func (a *SqlAdapter) GetUserByID(ctx context.Context, id uint64) (entity.User, error) {
	row, err := a.queries.GetUserByID(ctx, id)
	if err != nil {
		return entity.User{}, ports.ErrNotFound
	}
	return entity.User{
		ID:    row.ID,
		Name:  row.Name,
		Email: row.Email,
	}, nil
}

func (a *SqlAdapter) GetUserPlaylistIDs(ctx context.Context, userID uint64) ([]uint64, error) {
	return a.queries.GetUserPlaylistIDs(ctx, userID)
}
