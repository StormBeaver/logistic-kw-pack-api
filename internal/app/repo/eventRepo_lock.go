package repo

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/stormbeaver/logistic-pack-api/internal/model"
)

func (e eventRepo) Lock(ctx context.Context, count uint64) ([]model.PackEvent, error) {

	var (
		events = make([]RepoPackEvent, 0, count)
		ids    = make([]uint64, 0, count)
	)

	tx, err := e.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	err = AcquireLockTx(ctx, tx)

	if err != nil {
		return nil, fmt.Errorf("try Lock: %w", err)
	}

	sQuery := sq.Select("id", "type", "lock", "payload").
		From("packs_events").
		Where(sq.Eq{"lock": false}).
		Limit(count).
		RunWith(tx).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := sQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("convert to sql: %w", err)
	}

	err = tx.SelectContext(ctx, &events, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query in Lock: %w", err)
	}

	for _, v := range events {
		ids = append(ids, v.ID)
	}

	uQuery := sq.Update("packs_events").
		Set("lock", true).
		Set("updated", time.Now()).
		Where(sq.Eq{"id": ids}).
		RunWith(tx).
		PlaceholderFormat(sq.Dollar)

	uQuery.QueryRowContext(ctx)

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("tx.Commit: %w", err)
	}

	parsedEvents, err := parsePackEvent(events)
	if err != nil {
		return nil, fmt.Errorf("parse events: %w", err)
	}

	return parsedEvents, nil
}
