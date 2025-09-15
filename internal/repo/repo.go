package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/stormbeaver/logistic-pack-api/internal/model"
)

// Repo is DAO for Pack
type Repo interface {
	Add(ctx context.Context, name string) (uint64, error)
	Get(ctx context.Context, packID uint64) (*model.Pack, error)
	List(ctx context.Context, cursor uint64, limit uint64) ([]*model.Pack, error)
	Update(ctx context.Context, packId uint64, name string) (bool, error)
	Remove(ctx context.Context, packID uint64) (bool, error)
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}
