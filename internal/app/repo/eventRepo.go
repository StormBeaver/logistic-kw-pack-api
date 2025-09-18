package repo

import (
	"context"
	"encoding/json"
	"fmt"

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

type RepoPackEvent struct {
	ID     uint64  `db:"id"`
	Type   string  `db:"type"`
	Status bool    `db:"lock"`
	Entity []uint8 `db:"payload"`
}

func NewEventRepo(db *sqlx.DB) EventRepo {
	return &eventRepo{db: db}
}

// parse []RepoPackEvent into []model.PackEvent
func parsePackEvent(src []RepoPackEvent) ([]model.PackEvent, error) { //TODO: delete this and change Lock+PreProcess signature to use RepoPackEvent instead PackEvent and maybe delete status string
	parsedEvents := make([]model.PackEvent, len(src))
	for i, v := range src {
		parsedEvents[i].ID = v.ID
		parsedEvents[i].Status = v.Status
		parsedEvents[i].Type = v.Type
		if err := json.Unmarshal(v.Entity, &parsedEvents[i].Entity); err != nil {
			return nil, fmt.Errorf("unmarshal PeroPackEvent to model.PackEvent: %w", err)
		}
	}
	return parsedEvents, nil
}
