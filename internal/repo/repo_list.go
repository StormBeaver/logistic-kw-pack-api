package repo

import (
	"context"

	"github.com/stormbeaver/logistic-pack-api/internal/model"
)

func (r *repo) ListPacks(ctx context.Context) ([]*model.Pack, error) {
	return []*model.Pack{}, nil
}
