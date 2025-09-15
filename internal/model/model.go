package model

import (
	"time"
)

type Pack struct {
	ID      uint64    `db:"id"`
	Name    string    `db:"name"`
	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
}

type PackEvent struct {
	ID     uint64 `db:"id"`
	Type   string `db:"type"`
	Status bool   `db:"lock"`
	Entity *Pack  `db:"payload"`
}
