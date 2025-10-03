package repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/StormBeaver/logistic-pack-api/internal/model"
)

func (e eventRepo) PreProcess(ctx context.Context, count uint64) ([]model.PackEvent, error) {

	err := AcquireLock(ctx, e.db)

	if err != nil {
		return nil, fmt.Errorf("try lock PreProcess: %w", err)
	}

	sQuery := sq.Select("id", "type", "lock", "payload").
		From("packs_events").
		Where(sq.Eq{"lock": true}).
		Limit(count).
		RunWith(e.db).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := sQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("convert to sql: %w", err)
	}

	events := make([]RepoPackEvent, 0, count)

	err = e.db.SelectContext(ctx, &events, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query in Lock: %w", err)
	}

	parsedEvents, err := parsePackEvent(events)
	if err != nil {
		return nil, fmt.Errorf("parse events: %w", err)
	}

	return parsedEvents, Unlock(ctx, e.db)
}
