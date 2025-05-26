package model

import "fmt"

type EventType uint8

type EventStatus uint8

const (
	Created EventType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

type Pack struct {
	ID  uint64 `db:"id"`
	Foo uint64 `db:"foo"`
}

func (p *Pack) String() string {
	return fmt.Sprintf("ID: %d\nFoo: %v\n", p.ID, p.Foo)
}

type PackEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Pack
}
