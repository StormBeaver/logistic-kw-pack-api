package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"route255/logistic-kw-pack-api/internal/model"
)

// Repo is DAO for Pack
type Repo interface {
	DescribePack(ctx context.Context, packID uint64) (*model.Pack, error)
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}

func (r *repo) DescribePack(ctx context.Context, packID uint64) (*model.Pack, error) {
	return nil, nil
}
