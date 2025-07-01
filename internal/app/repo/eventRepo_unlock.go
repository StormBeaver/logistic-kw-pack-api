package repo

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (e eventRepo) Unlock(ctx context.Context, eventIDs []uint64) error {
	Query := sq.Update("packs_events").
		Set("lock", false).
		Set("updated", time.Now()).
		Where(sq.Eq{"id": eventIDs}).
		RunWith(e.db).
		PlaceholderFormat(sq.Dollar)

	Query.QueryRowContext(ctx)
	return nil
}
