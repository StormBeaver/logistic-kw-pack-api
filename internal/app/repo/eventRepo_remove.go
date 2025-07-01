package repo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (e eventRepo) Remove(ctx context.Context, eventIDs []uint64) error {
	Query := sq.Delete("packs_events").
		Where(sq.Eq{"id": eventIDs}).
		RunWith(e.db).
		PlaceholderFormat(sq.Dollar)

	Query.QueryRowContext(ctx)
	return nil
}
