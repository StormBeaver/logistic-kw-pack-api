package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/stormbeaver/logistic-pack-api/internal/model"
)

// Repo is DAO for Pack
type Repo interface {
	DescribePack(ctx context.Context, packID uint64) (*model.Pack, error)
	CreatePack(ctx context.Context, name string) (*model.Pack, error)
	RemovePack(ctx context.Context, packID uint64) (*model.Pack, error)
	ListPacks(ctx context.Context) ([]*model.Pack, error)
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}
