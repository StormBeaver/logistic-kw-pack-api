package repo

import (
	"context"

	"github.com/stormbeaver/logistic-pack-api/internal/model"
)

func (r *repo) ListPacks(ctx context.Context) ([]*model.Pack, error) {
	return []*model.Pack{
		{ID: 4, Name: "empty4"},
		{ID: 5, Name: "empty5"},
		{ID: 6, Name: "empty6"},
	}, nil
}
