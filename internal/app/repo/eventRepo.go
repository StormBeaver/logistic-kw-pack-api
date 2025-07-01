package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/stormbeaver/logistic-pack-api/internal/model"
)

type EventRepo interface {
	PreProcess(ctx context.Context, n uint64) ([]model.PackEvent, error)

	Lock(ctx context.Context, n uint64) ([]model.PackEvent, error)
	Unlock(ctx context.Context, eventIDs []uint64) error

	Remove(ctx context.Context, eventIDs []uint64) error
}

type eventRepo struct {
	db *sqlx.DB
}

func NewEventRepo(db *sqlx.DB) EventRepo {
	return &eventRepo{db: db}
}
