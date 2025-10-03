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

func (r *repo) Remove(ctx context.Context, packID uint64) (bool, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return false, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	pack := model.Pack{ID: packID, Updated: time.Now()}

	packQuery := sq.Update("packs").
		Set("removed", true).
		Set("updated", pack.Updated).
		Where(sq.Eq{"id": pack.ID, "removed": false}).
		Suffix("RETURNING name, created").
		RunWith(tx).
		PlaceholderFormat(sq.Dollar)

	if err := packQuery.QueryRowContext(ctx).Scan(&pack.Name, &pack.Created); err != nil {
		return false, fmt.Errorf("exec query: %w", err)
	}

	payload, err := json.Marshal(pack)
	if err != nil {
		return false, fmt.Errorf("json.marshal: %w", err)
	}

	eventQuery := sq.Insert("packs_events").
		Columns("pack_id", "type", "payload").
		Values(pack.ID, "removed", payload).
		RunWith(tx).
		PlaceholderFormat(sq.Dollar)

	eventQuery.QueryRowContext(ctx)

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("tx commit: %w", err)
	}

	return true, nil
}
