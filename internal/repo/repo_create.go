package repo

import (
	"context"

	"github.com/stormbeaver/logistic-pack-api/internal/model"
)

func (r *repo) CreatePack(ctx context.Context, name string) (*model.Pack, error) {
	return &model.Pack{}, nil
}
