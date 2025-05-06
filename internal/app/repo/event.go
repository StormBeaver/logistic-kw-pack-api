package repo

import (
	"route255/L-KW-P-API/internal/model"
)

type EventRepo interface {
	PreProcess(n uint64) ([]model.PackEvent, error)

	Lock(n uint64) ([]model.PackEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []model.PackEvent) error
	Remove(eventIDs []uint64) error
}
