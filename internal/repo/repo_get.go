package repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/StormBeaver/logistic-pack-api/internal/model"
)

func (r *repo) Get(ctx context.Context, packID uint64) (*model.Pack, error) {
	packQuery := sq.Select("id", "name", "created").
		From("packs").
		Where(sq.Eq{"id": packID, "removed": false}).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	var pack model.Pack
	if err := packQuery.QueryRowContext(ctx).Scan(&pack.ID, &pack.Name, &pack.Created); err != nil {
		return nil, fmt.Errorf("exec query and scan to pack: %w", err)
	}

	return &pack, nil
}
