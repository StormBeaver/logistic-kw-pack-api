package repo

import (
	"github.com/stormbeaver/logistic-pack-api/internal/model"
)

type EventRepo interface {
	PreProcess(n uint64) ([]model.PackEvent, error)

	Lock(n uint64) ([]model.PackEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []model.PackEvent) error
	Remove(eventIDs []uint64) error
}
