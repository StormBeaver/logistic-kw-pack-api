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
	Id       uint64
	Title    string
	Describe string
}

func (p *Pack) String() string {
	return fmt.Sprintf("ID: %d\nTitle: %s\nDescribe: %s", p.Id, p.Title, p.Describe)
}

type PackEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Pack
}
