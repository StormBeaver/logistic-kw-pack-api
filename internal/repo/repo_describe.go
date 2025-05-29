package repo

import (
	"context"

	"github.com/stormbeaver/logistic-pack-api/internal/model"
)

func (r *repo) DescribePack(ctx context.Context, packID uint64) (*model.Pack, error) {
	return &model.Pack{ID: 1, Name: "empty1"}, nil
}
