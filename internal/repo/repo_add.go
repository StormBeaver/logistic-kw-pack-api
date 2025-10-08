package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/StormBeaver/logistic-pack-api/internal/model"
)

// add entities into packs and packs_events tables
func (r *repo) Add(ctx context.Context, name, describe string) (uint64, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	pack := model.Pack{
		Name:     name,
		Describe: describe,
		Created:  time.Now(),
	}

	packQuery := sq.Insert("packs").
		Columns("name", "describe", "created").
		Values(pack.Name, pack.Describe, pack.Created).
		Suffix("RETURNING id").
		RunWith(tx).
		PlaceholderFormat(sq.Dollar)

	if err := packQuery.QueryRowContext(ctx).Scan(&pack.ID); err != nil {
		return 0, fmt.Errorf("exec query: %w", err)
	}

	payload, err := json.Marshal(pack)
	if err != nil {
		return 0, fmt.Errorf("json.marshal(pack): %w", err)
	}

	eventQuery := sq.Insert("packs_events").
		Columns("pack_id", "type", "payload").
		Values(pack.ID, "created", payload).
		RunWith(tx).
		PlaceholderFormat(sq.Dollar)

	eventQuery.QueryRowContext(ctx)

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("tx commit: %w", err)
	}

	return pack.ID, nil
}
