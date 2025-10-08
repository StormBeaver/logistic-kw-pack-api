package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/StormBeaver/logistic-pack-api/internal/model"
)

var errArgsEmpty = errors.New("arguments are empty")

// add entities into packs and packs_events tables
func (r *repo) Update(ctx context.Context, packID uint64, name, describe string) (bool, error) {

	if name == "" && describe == "" {
		return false, errArgsEmpty
	}

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return false, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	pack := model.Pack{
		ID:      packID,
		Updated: time.Now(),
	}

	updateQuery := sq.Update("packs").
		Set("updated", pack.Updated).
		Where(sq.Eq{"id": pack.ID, "removed": false}).
		Suffix("RETURNING name, describe, created").
		RunWith(tx).
		PlaceholderFormat(sq.Dollar)

	if describe != "" {
		updateQuery = updateQuery.Set("describe", describe)
	}

	if name != "" {
		updateQuery = updateQuery.Set("name", name)
	}

	if err := updateQuery.QueryRowContext(ctx).Scan(&pack.Name, &pack.Describe, &pack.Created); err != nil {
		return false, fmt.Errorf("exec query: %w", err)
	}

	payload, err := json.Marshal(pack)
	if err != nil {
		return false, fmt.Errorf("json.marshal(pack): %w", err)
	}

	eventQuery := sq.Insert("packs_events").
		Columns("pack_id", "type", "payload").
		Values(pack.ID, "updated", payload).
		RunWith(tx).
		PlaceholderFormat(sq.Dollar)

	eventQuery.QueryRowContext(ctx)

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("tx commit: %w", err)
	}

	return true, nil
}
