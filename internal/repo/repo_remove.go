package repo

import (
	"context"

	"github.com/stormbeaver/logistic-pack-api/internal/model"
)

func (r *repo) RemovePack(ctx context.Context, packID uint64) (*model.Pack, error) {
	return &model.Pack{ID: 3, Name: "empty3"}, nil
}
