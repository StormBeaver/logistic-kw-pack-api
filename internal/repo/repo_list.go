package repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/stormbeaver/logistic-pack-api/internal/model"
)

func (r *repo) List(ctx context.Context, cursor uint64, limit uint64) ([]*model.Pack, error) {
	packQuery := sq.Select("id", "name", "created").
		From("packs").
		Where(sq.Eq{"removed": false}).
		Offset(cursor - 1).
		Limit(limit).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := packQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("convert to sql: %w", err)
	}

	packs := make([]*model.Pack, 0, limit)

	err = r.db.SelectContext(ctx, &packs, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	return packs, nil
}
